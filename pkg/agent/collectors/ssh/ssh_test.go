package ssh

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// Mock data
const mockKeyRSA = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCQh08WbnI8YIATA/frSJGvW2oUjvZ387QDwH5OPfw8gN+SGpcWYdXBMVljVh72zPjiEBkeU9MXcEccKMGQXCE6AG5xsxA9saXJ3JIB/Ydo1wgabWXw3/HbwE41PWNu0l4xcM3d2IJRMxqk1HB5ipwvWkbZ9Xs3XqngYIblHKkdkCAq4cZ6EAWA2LgvNcUZh31IN0m3d6ZcOF4Xgg/qCIEoiPXM2GDRttASXjLNpEiB1sF/lbKCEyzmrejNUMEOiw/Hlwid+y5Vmmefn7sKVoCZ42ZdDuynzFWR+fs4ISjMoJvPKrtFoea+JMdO9hNp4QUuwbYDh8CJNTq1pV/8UTJOFk4FShrncYnWlgeRmKp+P39QE/JFJMwxErkGT/mTH38l9QYIgeflNpsOTcM3fpprSxU/sVIxb3IFssQQwNQQ9Bp+1eo2nwrjxuY1EyZiqfUsAo+OdRpcQwV+2LyNneNlp7Az3n5/xuYgwTwKUslPdhlA4h/mjfdF9EAIl20CwZM= user@host`

func TestGetSshKeyInfo_ValidKey(t *testing.T) {
	keyList := []string{mockKeyRSA}
	keyPaths := []string{"/home/user/.ssh/id_rsa.pub"}
	logger := &logrus.Logger{}
	result := GetSshKeyInfo(logger, keyList, keyPaths)

	assert.Len(t, result, 1)
	assert.Equal(t, "user@host", result[0].Comment)
	assert.Equal(t, "ssh-rsa", result[0].Type)
	assert.Equal(t, "/home/user", result[0].Path)
	assert.Equal(t, "id_rsa.pub", result[0].Name)
	assert.False(t, result[0].FromAuthorizedKeys)
}

func TestGetSshKeyInfo_InvalidKey(t *testing.T) {
	keyList := []string{"invalid-key-data"}
	keyPaths := []string{"/home/user/.ssh/bad.pub"}
	logger := &logrus.Logger{}
	result := GetSshKeyInfo(logger, keyList, keyPaths)

	assert.Len(t, result, 0)
}

func TestGetPath(t *testing.T) {
	base, file := GetPath("/home/user/.ssh/id_rsa.pub")

	assert.Equal(t, "/home/user", base)
	assert.Equal(t, "id_rsa.pub", file)
}

func TestContains(t *testing.T) {
	values := []string{"one", "two", "three"}

	assert.True(t, Contains(values, "two"))
	assert.False(t, Contains(values, "four"))
}

func TestIsPrivateKeyFile_RSA(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "id_rsa")
	content := "-----BEGIN RSA PRIVATE KEY-----\nMIIBOgIBAAJBALw=\n-----END RSA PRIVATE KEY-----"
	err := os.WriteFile(file, []byte(content), 0600)
	assert.NoError(t, err)
	assert.True(t, isPrivateKeyFile(file))
}

func TestIsPrivateKeyFile_OpenSSH(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "id_ed25519")
	content := "-----BEGIN OPENSSH PRIVATE KEY-----\nabcde\n-----END OPENSSH PRIVATE KEY-----"
	err := os.WriteFile(file, []byte(content), 0600)
	assert.NoError(t, err)
	assert.True(t, isPrivateKeyFile(file))
}

func TestIsPrivateKeyFile_NotPrivateKey(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "id_rsa.pub")
	content := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7"
	err := os.WriteFile(file, []byte(content), 0600)
	assert.NoError(t, err)
	assert.False(t, isPrivateKeyFile(file))
}

func TestIsPrivateKeyFile_FileDoesNotExist(t *testing.T) {
	file := "/nonexistent/path/to/key"
	assert.False(t, isPrivateKeyFile(file))
}

func TestIsPrivateKeyFile_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "empty")
	err := os.WriteFile(file, []byte(""), 0600)
	assert.NoError(t, err)
	assert.False(t, isPrivateKeyFile(file))
}
func TestParseFile_PubAndAuthorizedKeys(t *testing.T) {
	dir := t.TempDir()

	// Create a .pub file
	pubFile := filepath.Join(dir, "id_rsa.pub")
	err := os.WriteFile(pubFile, []byte("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7"), 0644)
	assert.NoError(t, err)

	// Create an authorized_keys file
	authFile := filepath.Join(dir, "authorized_keys")
	err = os.WriteFile(authFile, []byte("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7 user@host"), 0644)
	assert.NoError(t, err)

	// Create a non-regular file (directory)
	nonRegular := filepath.Join(dir, "not_a_file")
	err = os.Mkdir(nonRegular, 0755)
	assert.NoError(t, err)

	entries, err := os.ReadDir(dir)
	assert.NoError(t, err)

	logger := &logrus.Logger{}
	keys, knownHosts := parseFile(logger, entries, dir)

	assert.Contains(t, keys, pubFile)
	assert.Contains(t, keys, authFile)
	assert.Len(t, knownHosts, 0)
}

func TestParseFile_PrivateKeySkipped(t *testing.T) {
	dir := t.TempDir()

	// Create a private key file
	privFile := filepath.Join(dir, "id_rsa")
	content := "-----BEGIN RSA PRIVATE KEY-----\nMIIBOgIBAAJBALw=\n-----END RSA PRIVATE KEY-----"
	err := os.WriteFile(privFile, []byte(content), 0600)
	assert.NoError(t, err)

	entries, err := os.ReadDir(dir)
	assert.NoError(t, err)

	logger := &logrus.Logger{}
	keys, knownHosts := parseFile(logger, entries, dir)

	assert.NotContains(t, keys, privFile)
	assert.Len(t, knownHosts, 0)
}

func TestParseFile_KnownHosts(t *testing.T) {
	dir := t.TempDir()

	// Create a known_hosts file with a valid entry
	knownHostsFile := filepath.Join(dir, "known_hosts")
	knownHostsContent := "host.example.com " + mockKeyRSA
	err := os.WriteFile(knownHostsFile, []byte(knownHostsContent), 0644)
	assert.NoError(t, err)

	entries, err := os.ReadDir(dir)
	assert.NoError(t, err)

	logger := &logrus.Logger{}
	keys, knownHosts := parseFile(logger, entries, dir)

	assert.Len(t, keys, 0)
	assert.Len(t, knownHosts, 1)
	assert.Contains(t, knownHosts[0].Hostname, "host.example.com")
	assert.Equal(t, "ssh-rsa", knownHosts[0].Type)
}

func TestParseFile_KnownHostsFail(t *testing.T) {
	dir := t.TempDir()

	// Create a known_hosts file with a valid entry
	knownHostsFile := filepath.Join(dir, "known_hosts")
	knownHostsContent := "host.example.com ssh-rsa aaa"
	err := os.WriteFile(knownHostsFile, []byte(knownHostsContent), 0644)
	assert.NoError(t, err)

	entries, err := os.ReadDir(dir)
	assert.NoError(t, err)

	logger := &logrus.Logger{}
	keys, knownHosts := parseFile(logger, entries, dir)

	assert.Len(t, keys, 0)
	assert.Len(t, knownHosts, 0)
}

func TestParseFile_SkipsNonRegularFiles(t *testing.T) {
	dir := t.TempDir()

	// Create a directory entry
	subDir := filepath.Join(dir, "subdir")
	err := os.Mkdir(subDir, 0755)
	assert.NoError(t, err)

	entries, err := os.ReadDir(dir)
	assert.NoError(t, err)

	logger := &logrus.Logger{}
	keys, knownHosts := parseFile(logger, entries, dir)

	assert.Len(t, keys, 0)
	assert.Len(t, knownHosts, 0)
}

func TestParseFile_UnknownFileType(t *testing.T) {
	dir := t.TempDir()

	// Create a file with unknown extension
	unknownFile := filepath.Join(dir, "randomfile.txt")
	err := os.WriteFile(unknownFile, []byte("some data"), 0644)
	assert.NoError(t, err)

	entries, err := os.ReadDir(dir)
	assert.NoError(t, err)

	logger := &logrus.Logger{}
	keys, knownHosts := parseFile(logger, entries, dir)

	assert.Len(t, keys, 0)
	assert.Len(t, knownHosts, 0)
}
func TestReadPubKeyFile_ValidPubKey(t *testing.T) {
	dir := t.TempDir()
	pubFile := filepath.Join(dir, "id_rsa.pub")
	content := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCQh08WbnI8YIATA/frSJGvW2oUjvZ387QDwH5OPfw8gN+SGpcWYdXBMVljVh72zPjiEBkeU9MXcEccKMGQXCE6AG5xsxA9saXJ3JIB/Ydo1wgabWXw3/HbwE41PWNu0l4xcM3d2IJRMxqk1HB5ipwvWkbZ9Xs3XqngYIblHKkdkCAq4cZ6EAWA2LgvNcUZh31IN0m3d6ZcOF4Xgg/qCIEoiPXM2GDRttASXjLNpEiB1sF/lbKCEyzmrejNUMEOiw/Hlwid+y5Vmmefn7sKVoCZ42ZdDuynzFWR+fs4ISjMoJvPKrtFoea+JMdO9hNp4QUuwbYDh8CJNTq1pV/8UTJOFk4FShrncYnWlgeRmKp+P39QE/JFJMwxErkGT/mTH38l9QYIgeflNpsOTcM3fpprSxU/sVIxb3IFssQQwNQQ9Bp+1eo2nwrjxuY1EyZiqfUsAo+OdRpcQwV+2LyNneNlp7Az3n5/xuYgwTwKUslPdhlA4h/mjfdF9EAIl20CwZM= user@host"
	err := os.WriteFile(pubFile, []byte(content), 0644)
	assert.NoError(t, err)

	logger := &logrus.Logger{}
	result := ReadPubKeyFile(logger, []string{pubFile})
	assert.Len(t, result, 1)
	assert.Equal(t, "user@host", result[0].Comment)
	assert.Equal(t, "ssh-rsa", result[0].Type)
	assert.Equal(t, dir, result[0].Path)
	assert.Equal(t, "id_rsa.pub", result[0].Name)
}

func TestReadPubKeyFile_InvalidPubKey(t *testing.T) {
	dir := t.TempDir()
	pubFile := filepath.Join(dir, "id_rsa.pub")
	content := "not-a-valid-key"
	err := os.WriteFile(pubFile, []byte(content), 0644)
	assert.NoError(t, err)

	logger := &logrus.Logger{}
	result := ReadPubKeyFile(logger, []string{pubFile})

	assert.Len(t, result, 0)
}

func TestReadPubKeyFile_FileDoesNotExist(t *testing.T) {
	logger := &logrus.Logger{}
	result := ReadPubKeyFile(logger, []string{"/nonexistent/path/to/key.pub"})
	assert.Len(t, result, 0)
}

func TestReadPubKeyFile_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	pubFile := filepath.Join(dir, "empty.pub")
	err := os.WriteFile(pubFile, []byte(""), 0644)
	assert.NoError(t, err)

	logger := &logrus.Logger{}
	result := ReadPubKeyFile(logger, []string{pubFile})

	assert.Len(t, result, 0)
}

func TestReadPubKeyFile_MultipleFiles(t *testing.T) {
	dir := t.TempDir()
	pubFile1 := filepath.Join(dir, "id_rsa.pub")
	pubFile2 := filepath.Join(dir, "id_ecdsa.pub")
	content1 := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCQh08WbnI8YIATA/frSJGvW2oUjvZ387QDwH5OPfw8gN+SGpcWYdXBMVljVh72zPjiEBkeU9MXcEccKMGQXCE6AG5xsxA9saXJ3JIB/Ydo1wgabWXw3/HbwE41PWNu0l4xcM3d2IJRMxqk1HB5ipwvWkbZ9Xs3XqngYIblHKkdkCAq4cZ6EAWA2LgvNcUZh31IN0m3d6ZcOF4Xgg/qCIEoiPXM2GDRttASXjLNpEiB1sF/lbKCEyzmrejNUMEOiw/Hlwid+y5Vmmefn7sKVoCZ42ZdDuynzFWR+fs4ISjMoJvPKrtFoea+JMdO9hNp4QUuwbYDh8CJNTq1pV/8UTJOFk4FShrncYnWlgeRmKp+P39QE/JFJMwxErkGT/mTH38l9QYIgeflNpsOTcM3fpprSxU/sVIxb3IFssQQwNQQ9Bp+1eo2nwrjxuY1EyZiqfUsAo+OdRpcQwV+2LyNneNlp7Az3n5/xuYgwTwKUslPdhlA4h/mjfdF9EAIl20CwZM= user@host"
	content2 := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAICvQw1tQw1tQw1tQw1tQw1tQw1tQw1tQw1tQw1tQw1tQ user2@host"
	err := os.WriteFile(pubFile1, []byte(content1), 0644)
	assert.NoError(t, err)
	err = os.WriteFile(pubFile2, []byte(content2), 0644)
	assert.NoError(t, err)

	logger := &logrus.Logger{}
	result := ReadPubKeyFile(logger, []string{pubFile1, pubFile2})

	assert.Len(t, result, 2)
	assert.Equal(t, "user@host", result[0].Comment)
	assert.Equal(t, "user2@host", result[1].Comment)
}
func TestGetSshFiles_MatchPubFiles(t *testing.T) {
	dir := t.TempDir()
	// Create files
	pubFile := filepath.Join(dir, "id_rsa.pub")
	privFile := filepath.Join(dir, "id_rsa")
	txtFile := filepath.Join(dir, "notes.txt")
	err := os.WriteFile(pubFile, []byte("pubkey"), 0644)
	assert.NoError(t, err)
	err = os.WriteFile(privFile, []byte("privkey"), 0600)
	assert.NoError(t, err)
	err = os.WriteFile(txtFile, []byte("text"), 0644)
	assert.NoError(t, err)

	// Match only .pub files
	matchFunc := func(name string) bool {
		return filepath.Ext(name) == ".pub"
	}

	files, err := GetSshFiles(dir, matchFunc)
	assert.NoError(t, err)
	assert.Len(t, files, 1)
	assert.Contains(t, files, pubFile)
}

func TestGetSshFiles_MatchAllFiles(t *testing.T) {
	dir := t.TempDir()
	file1 := filepath.Join(dir, "file1")
	file2 := filepath.Join(dir, "file2")
	err := os.WriteFile(file1, []byte("data1"), 0644)
	assert.NoError(t, err)
	err = os.WriteFile(file2, []byte("data2"), 0644)
	assert.NoError(t, err)

	matchFunc := func(name string) bool { return true }

	files, err := GetSshFiles(dir, matchFunc)
	assert.NoError(t, err)
	assert.Len(t, files, 2)
	assert.Contains(t, files, file1)
	assert.Contains(t, files, file2)
}

func TestGetSshFiles_NoMatch(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "id_rsa.pub")
	err := os.WriteFile(file, []byte("pubkey"), 0644)
	assert.NoError(t, err)

	matchFunc := func(name string) bool { return false }

	files, err := GetSshFiles(dir, matchFunc)
	assert.NoError(t, err)
	assert.Len(t, files, 0)
}

func TestGetSshFiles_SkipsNonRegularFiles(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "subdir")
	err := os.Mkdir(subDir, 0744)
	assert.NoError(t, err)
	file := filepath.Join(dir, "id_rsa.pub")
	err = os.WriteFile(file, []byte("pubkey"), 0644)
	assert.NoError(t, err)

	matchFunc := func(name string) bool { return true }

	files, err := GetSshFiles(dir, matchFunc)
	assert.NoError(t, err)
	assert.Len(t, files, 1)
	assert.Contains(t, files, file)
}

func TestGetSshFiles_DirDoesNotExist(t *testing.T) {
	nonExistentDir := filepath.Join(os.TempDir(), "does_not_exist_ssh_test")
	matchFunc := func(name string) bool { return true }

	files, err := GetSshFiles(nonExistentDir, matchFunc)
	assert.Error(t, err)
	assert.Nil(t, files)
}

func TestGetSshInfo_DirDoesNotExist(t *testing.T) {
	logger := logrus.New()
	var homeDirList []string
	homeDirList = append(homeDirList, "/home/fake")
	a, b, c := GetSshInfo(logger, homeDirList)
	assert.Empty(t, a)
	assert.Empty(t, b)
	assert.Empty(t, c)

}

func TestGetAllSshFiles_DirDoesNotExist(t *testing.T) {
	logger := logrus.New()
	var homeDirList []string
	homeDirList = append(homeDirList, "/home/fake")
	a, b := GetAllSshFiles(logger, homeDirList)
	assert.Empty(t, a)
	assert.Empty(t, b)

}

func TestGetAllSshFiles(t *testing.T) {
	dir := t.TempDir()
	dirSsh := dir + "/.ssh"
	err := os.Mkdir(dirSsh, 0744)
	assert.NoError(t, err)
	pubFile1 := filepath.Join(dirSsh, "id_rsa.pub")
	content1 := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCQh08WbnI8YIATA/frSJGvW2oUjvZ387QDwH5OPfw8gN+SGpcWYdXBMVljVh72zPjiEBkeU9MXcEccKMGQXCE6AG5xsxA9saXJ3JIB/Ydo1wgabWXw3/HbwE41PWNu0l4xcM3d2IJRMxqk1HB5ipwvWkbZ9Xs3XqngYIblHKkdkCAq4cZ6EAWA2LgvNcUZh31IN0m3d6ZcOF4Xgg/qCIEoiPXM2GDRttASXjLNpEiB1sF/lbKCEyzmrejNUMEOiw/Hlwid+y5Vmmefn7sKVoCZ42ZdDuynzFWR+fs4ISjMoJvPKrtFoea+JMdO9hNp4QUuwbYDh8CJNTq1pV/8UTJOFk4FShrncYnWlgeRmKp+P39QE/JFJMwxErkGT/mTH38l9QYIgeflNpsOTcM3fpprSxU/sVIxb3IFssQQwNQQ9Bp+1eo2nwrjxuY1EyZiqfUsAo+OdRpcQwV+2LyNneNlp7Az3n5/xuYgwTwKUslPdhlA4h/mjfdF9EAIl20CwZM= user@host"
	err = os.WriteFile(pubFile1, []byte(content1), 0644)
	assert.NoError(t, err)

	knownHostsFile := filepath.Join(dirSsh, "known_hosts")
	knownHostsContent := "host.example.com " + mockKeyRSA
	err = os.WriteFile(knownHostsFile, []byte(knownHostsContent), 0644)
	assert.NoError(t, err)

	logger := logrus.New()
	var homeDirList []string
	homeDirList = append(homeDirList, dir)
	a, b := GetAllSshFiles(logger, homeDirList)
	assert.Equal(t, a[0], pubFile1)
	assert.Equal(t, b[0].Hostname, "host.example.com")
	assert.Equal(t, b[0].Type, "ssh-rsa")
	assert.Equal(t, b[0].Fingerprint, "SHA256:onumzaYU98XSgnDrr1a2yBjuArXo+EbG/MMRA/bhtG8")

}

func TestParseAuthFiles_DoesNotExist(t *testing.T) {
	logger := logrus.New()
	var authFile []string
	authFile = append(authFile, "/home/fake/.ssh/authorized_keys")
	authFile = append(authFile, "/tmp/test")

	a := ParseAuthFiles(logger, authFile)
	assert.Nil(t, a)
}

func TestParseAuthFiles(t *testing.T) {
	dir := t.TempDir()
	dirSsh := dir + "/.ssh"
	err := os.Mkdir(dirSsh, 0744)
	assert.NoError(t, err)
	pubFile1 := filepath.Join(dirSsh, "authorized_keys")
	content1 := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCQh08WbnI8YIATA/frSJGvW2oUjvZ387QDwH5OPfw8gN+SGpcWYdXBMVljVh72zPjiEBkeU9MXcEccKMGQXCE6AG5xsxA9saXJ3JIB/Ydo1wgabWXw3/HbwE41PWNu0l4xcM3d2IJRMxqk1HB5ipwvWkbZ9Xs3XqngYIblHKkdkCAq4cZ6EAWA2LgvNcUZh31IN0m3d6ZcOF4Xgg/qCIEoiPXM2GDRttASXjLNpEiB1sF/lbKCEyzmrejNUMEOiw/Hlwid+y5Vmmefn7sKVoCZ42ZdDuynzFWR+fs4ISjMoJvPKrtFoea+JMdO9hNp4QUuwbYDh8CJNTq1pV/8UTJOFk4FShrncYnWlgeRmKp+P39QE/JFJMwxErkGT/mTH38l9QYIgeflNpsOTcM3fpprSxU/sVIxb3IFssQQwNQQ9Bp+1eo2nwrjxuY1EyZiqfUsAo+OdRpcQwV+2LyNneNlp7Az3n5/xuYgwTwKUslPdhlA4h/mjfdF9EAIl20CwZM= user@host"
	err = os.WriteFile(pubFile1, []byte(content1), 0644)
	assert.NoError(t, err)
	logger := logrus.New()
	var authFile []string
	authFile = append(authFile, pubFile1)
	authFile = append(authFile, "/tmp/test")

	a := ParseAuthFiles(logger, authFile)
	assert.Len(t, a, 1)
}
