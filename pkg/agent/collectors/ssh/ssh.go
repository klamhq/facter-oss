package ssh

import (
	"bufio"
	"bytes"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"io"
	"os"
	"path/filepath"
	"strings"

	"slices"

	"github.com/klamhq/facter-oss/pkg/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

// GetPath from ssh keys or autorized keys files
func GetPath(path string) (string, string) {
	basePath, fileName := filepath.Split(path)
	pathSplit := strings.Split(string(basePath), ".ssh")
	pathSplitSuffix := strings.TrimSuffix(pathSplit[0], "/")
	return pathSplitSuffix, fileName
}

// GetSshKeyInfo retrieves SSH key information from provided key lists and paths.
// It unmarshals the public keys, extracts their properties, and returns a slice of SshKeyInfo models.
// It ensures that each key is processed only once, avoiding duplicates based on the key value and path.
func GetSshKeyInfo(logger *logrus.Logger, keyLists, keyPaths []string) []models.SshKeyInfo {
	var sshKeys []models.SshKeyInfo
	var seenKeys = make(map[string]struct{})
	var seenPaths = make(map[string]struct{})

	for _, key := range keyLists {
		for _, path := range keyPaths {
			if _, ok := seenKeys[key]; ok {
				continue
			}
			if _, ok := seenPaths[path]; ok {
				continue
			}

			fromAuthorizedKeys := false
			basePath, fileName := GetPath(path)
			if fileName == "authorized_keys" {
				fromAuthorizedKeys = true
			}
			Owner := getUsernameFromPath(path)
			pubKey, comment, options, _, err := ssh.ParseAuthorizedKey([]byte(key))
			if err != nil {
				logger.Warnf("Unable to parse key %s: %v", key, err)
				continue
			}

			parsedCryptoKey := pubKey.(ssh.CryptoPublicKey).CryptoPublicKey()
			var length int
			switch k := parsedCryptoKey.(type) {
			case *rsa.PublicKey:
				length = k.Size() * 8 // La taille de la clé RSA est en octets, donc on multiplie par 8 pour obtenir les bits
			case *ecdsa.PublicKey:
				length = k.Params().BitSize
			case ed25519.PublicKey:
				// La taille n'est pas mesurée en "bits de clé" pour Ed25519, elle est toujours 256 bits
				length = 256
			default:
				length = 0
			}

			info := models.SshKeyInfo{
				Comment:            comment,
				Length:             int64(length),
				Type:               pubKey.Type(),
				Fingerprint:        ssh.FingerprintSHA256(pubKey),
				Options:            options,
				Path:               basePath,
				Name:               fileName,
				FromAuthorizedKeys: fromAuthorizedKeys,
				Owner:              Owner,
			}

			if info.Name != "" && info.Type != "" {
				sshKeys = append(sshKeys, info)
				seenKeys[key] = struct{}{}
				seenPaths[path] = struct{}{}
			}
		}
	}

	return sshKeys
}

// Contains checks if a string is present in a slice of strings.
func Contains(homeDir []string, str string) bool {
	return slices.Contains(homeDir, str)
}

// isPrivateKeyFile checks if the given file path is a private key file.
func isPrivateKeyFile(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		// On considère que si on ne peut pas lire, on ignore par sécurité
		return false
	}
	defer f.Close()

	// Lire les premiers 1 Ko max (ça suffit pour voir l’en-tête PEM)
	buf := make([]byte, 1024)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return false
	}

	content := string(buf[:n])

	// Vérifier les en-têtes connus de clés privées
	privateKeyHeaders := []string{
		"-----BEGIN OPENSSH PRIVATE KEY-----",
		"-----BEGIN RSA PRIVATE KEY-----",
		"-----BEGIN DSA PRIVATE KEY-----",
		"-----BEGIN EC PRIVATE KEY-----",
		"-----BEGIN ED25519 PRIVATE KEY-----",
	}
	for _, header := range privateKeyHeaders {
		if strings.Contains(content, header) {
			return true
		}
	}
	return false
}

// parseFile parses the directory entries for SSH key files.
// It looks for files with .pub extension, known_host or an authorized_keys file.
func parseFile(logger *logrus.Logger, entries []os.DirEntry, dirname string) ([]string, []models.KnownHost) {
	var keyPub []string
	var knownHost []models.KnownHost
	for _, entry := range entries {
		if entry.Type().IsRegular() {
			name := entry.Name()
			fullPath := filepath.Join(dirname, name)

			if isPrivateKeyFile(fullPath) {
				logger.Debugf("Skipping private key file: %v", fullPath)
				continue
			}
			switch filepath.Ext(name) {
			case ".pub":
				keyPub = append(keyPub, fullPath)
			default:
				switch name {
				case "authorized_keys":
					logger.Debugf("Found authorized_keys file: %v", fullPath)
					keyPub = append(keyPub, fullPath)
				case "known_hosts":
					knownHost = parseKnownHostsSkipHashed(logger, entry, dirname)
				default:
					logger.Debugf("File %v does not match expected SSH key file patterns", fullPath)
				}
			}
		} else {
			logger.Debugf("Skipping non-regular file: %v", entry.Name())
		}
	}
	return keyPub, knownHost
}

// parseKnownHostsSkipHashed parses known_hosts files, skipping hashed entries.
func parseKnownHostsSkipHashed(logger *logrus.Logger, entry os.DirEntry, dirname string) []models.KnownHost {
	fullPath := filepath.Join(dirname, entry.Name())
	Owner := getUsernameFromPath(fullPath)
	file, err := os.Open(fullPath)
	if err != nil {
		logger.Errorf("Failed to open file %s: %v", fullPath, err)
		return nil
	}
	defer file.Close()

	var knownHost []models.KnownHost
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		_, hosts, pubKey, _, _, err := ssh.ParseKnownHosts([]byte(line))
		if err != nil {
			logger.Warnf("Error parsing line in %s: %v", fullPath, err)
			continue
		}

		// Vérifier si la ligne est hashée : on la skippe
		if len(hosts) == 1 && strings.HasPrefix(hosts[0], "|1|") {
			logger.Debugf("Skipping hashed known_hosts entry: %v", hosts[0])
			continue
		}
		knownHost = append(knownHost, models.KnownHost{
			Hostname:    strings.Join(hosts, ", "),
			Type:        pubKey.Type(),
			Fingerprint: ssh.FingerprintSHA256(pubKey),
			Owner:       Owner,
		})
	}

	if err := scanner.Err(); err != nil {
		logger.Errorf("Error reading file %s: %v", fullPath, err)
	}
	return knownHost
}

// ReadPubKeyFile reads public key files from the specified paths.
// It returns a slice of SshKeyInfo models containing the key values and their paths.
func ReadPubKeyFile(logger *logrus.Logger, pubKeyFiles []string) []models.SshKeyInfo {
	var keyValues []string
	var keyPaths []string

	for _, file := range pubKeyFiles {
		logger.Debugf("Reading public key file: %s", file)
		content, err := os.ReadFile(file)
		if err != nil {
			logger.Warnf("Cannot read public key file %s: %v", file, err)
			continue
		}
		keyPaths = append(keyPaths, file)
		keyValues = append(keyValues, string(content))
	}

	return GetSshKeyInfo(logger, keyValues, keyPaths)
}

// GetSshFiles retrieves SSH files from the specified directory based on a matching function.
// It returns a slice of file paths that match the criteria defined by the matchFunc.
func GetSshFiles(dir string, matchFunc func(name string) bool) ([]string, error) {
	var matches []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.Type().IsRegular() && matchFunc(entry.Name()) {
			matches = append(matches, filepath.Join(dir, entry.Name()))
		}
	}

	return matches, nil
}

// GetSshInfo retrieves SSH key information from the specified home directories.
// It reads public key files and authorized keys files, returning slices of SshKeyInfo and SshKeyAccess models.
func GetSshInfo(logger *logrus.Logger, homeDirs []string) ([]models.SshKeyInfo, []models.SshKeyAccess, []models.KnownHost) {
	logger.Debug("Fetching SSH key information from home directories")

	files, knownHost := GetAllSshFiles(logger, homeDirs)

	pubKeyInfos := ReadPubKeyFile(logger, files)

	keyAccess := ParseAuthFiles(logger, files)

	return pubKeyInfos, keyAccess, knownHost
}

// GetAllSshFiles retrieves all SSH key files from the specified home directories.
// It looks for files with .pub extension or an authorized_keys file in each .ssh directory.
func GetAllSshFiles(logger *logrus.Logger, homeDirs []string) ([]string, []models.KnownHost) {
	var allKeys []string
	var knownHost []models.KnownHost
	var parsed []string
	for _, homeDir := range homeDirs {
		sshDir := filepath.Join(homeDir, ".ssh")

		var files []os.DirEntry
		if _, err := os.Stat(sshDir); !os.IsNotExist(err) {
			logger.Debugf("Checking .ssh directory: %s", sshDir)
			files, err = os.ReadDir(sshDir)
			if err != nil {
				logger.Debugf("Unable to read .ssh directory %s: %v", sshDir, err)
				continue
			}
		}

		parsed, knownHost = parseFile(logger, files, sshDir)
		allKeys = append(allKeys, parsed...)
	}

	return allKeys, knownHost
}

// ParseAuthFiles reads authorized keys files and extracts SSH key access information.
// It returns a slice of SshKeyAccess models containing the fingerprint and associated user.
func ParseAuthFiles(logger *logrus.Logger, authFiles []string) []models.SshKeyAccess {
	var sshAuthKey []models.SshKeyAccess

	for _, authorizedKeyFile := range authFiles {
		authFile, err := os.ReadFile(authorizedKeyFile)
		if err != nil {
			logger.Warnf("Unable to read authorized keys file %s: %v", authorizedKeyFile, err)
			continue
		}

		scanner := bufio.NewScanner(bytes.NewReader(authFile))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue // ignore empty lines and comments
			}

			pk, _, _, _, err := ssh.ParseAuthorizedKey([]byte(line))
			if err != nil {
				logger.Warnf("Unable to parse authorized key %s in file %s: %v", line, authorizedKeyFile, err)
				continue
			}

			asUser := getUsernameFromPath(authorizedKeyFile)
			u := models.SshKeyAccess{
				Fingerprint: ssh.FingerprintLegacyMD5(pk),
				AsUser:      asUser,
			}
			sshAuthKey = append(sshAuthKey, u)
		}

		if err := scanner.Err(); err != nil {
			logger.Warnf("Error reading authorized keys file %s: %v", authorizedKeyFile, err)
		}
	}

	return sshAuthKey
}

// getUsernameFromPath extracts the username from a given file path.
// It looks for the "home" or "root" directory in the path and returns the username.
func getUsernameFromPath(path string) string {
	cleanPath := filepath.Clean(path)
	parts := strings.Split(cleanPath, string(os.PathSeparator))
	for i := 0; i < len(parts); i++ {
		switch parts[i] {
		case "home":
			if i+1 < len(parts) {
				return parts[i+1]
			}
		case "root":
			return "root"
		}
	}
	return ""
}
