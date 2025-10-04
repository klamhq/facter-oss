package users

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"slices"

	"github.com/klamhq/facter-oss/pkg/models"
	"github.com/klamhq/facter-oss/pkg/utils"
	"github.com/shirou/gopsutil/host"
	"github.com/sirupsen/logrus"
)

// GetConnectedUsers fetches the currently connected users on the system.
// It uses the gopsutil library to retrieve user statistics, which includes information about
// the terminal, host, and start time of each user's session.
// The function returns a slice of host.UserStat, which contains the details of each connected user.
// If an error occurs while fetching the user statistics, it logs the error and returns an empty slice.
func GetConnectedUsers(logger *logrus.Logger) []host.UserStat {
	ctx := context.Background()
	users, err := host.UsersWithContext(ctx)
	if err != nil {
		logger.WithError(err).Warnf("unable to fetch connected users: %v", err)
	}
	return users
}

// MergeUsersAndSessions merges system users with their connected sessions.
// It takes a slice of system users and a slice of connected user sessions,
// and returns a slice of system users with their sessions populated.
// If a system user has no connected sessions, it will have a session with "none" values.
// The function assumes that the connected users are represented as host.UserStat,
// which contains information about the terminal, host, and start time of the session.
// Each system user will have a Session field populated with the connected sessions or a default session if none are found.
func MergeUsersAndSessions(systemUsers []models.User, connectedUsers []host.UserStat) []models.User {

	sessionsMap := make(map[string][]host.UserStat)
	for _, cu := range connectedUsers {
		sessionsMap[cu.User] = append(sessionsMap[cu.User], cu)
	}

	// Pour chaque utilisateur système, on ajoute ses sessions connectées
	for i, su := range systemUsers {
		if sessions, ok := sessionsMap[su.Username]; ok {
			var convertedSessions []models.Session
			for _, s := range sessions { // sessions est []host.UserStat
				convertedSessions = append(convertedSessions, models.Session{
					Terminal:  s.Terminal,
					Host:      s.Host,
					Started:   int64(s.Started),
					Connected: true,
				})
			}
			systemUsers[i].Session = convertedSessions
		} else {
			var session []models.Session
			session = append(session, models.Session{
				Terminal:  "none",
				Host:      "none",
				Started:   0,
				Connected: false,
			})
			systemUsers[i].Session = session
		}
	}

	return systemUsers
}

// GetSystemUsers fetches the system users from the /etc/passwd file.
// It reads the file, parses each line to create a User struct, and returns a slice of User structs.
// If an error occurs while reading the file or parsing a line, it returns an empty slice and an error message.
func GetSystemUsers(passwdFilename string, logger *logrus.Logger) ([]models.User, error) {
	users, err := getUsers(passwdFilename, logger)
	if err != nil {
		return make([]models.User, 0), fmt.Errorf("unable to fetch users : %v", err)
	}

	return users, nil
}

// getUsers reads the /etc/passwd file and returns a slice of User structs.
// It parses each line of the file to extract user information such as username, name, GID, home directory, and UID.
// If an error occurs while reading the file or parsing a line, it logs a warning and continues to the next line.
func getUsers(passwdFilename string, logger *logrus.Logger) ([]models.User, error) {
	users := new([]models.User)
	passwdFile, err := os.ReadFile(passwdFilename)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch user list :%s", err)
	}
	for line := range strings.SplitSeq(string(passwdFile), "\n") {
		newUser, err := parseUserLine(line, logger)
		if newUser == nil {
			// Skip the first empty line
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, "#") {
				continue
			}
			if err != nil {
				logger.Warnf("error while reading user %v", err)
				continue
			}
		}
		*users = append(*users, *newUser)
	}
	return *users, nil
}

// parseUserLine convertit une ligne du fichier passwd en struct User.
func parseUserLine(line string, logger *logrus.Logger) (*models.User, error) {
	parts := strings.Split(line, ":")
	if len(parts) < 7 {
		return nil, fmt.Errorf("invalid line: %s", line)
	}
	username := parts[0]
	if username == "" {
		return nil, nil
	}

	uid := parts[2]
	gid := parts[3]
	home := parts[5]

	// Récupération du shell
	shell := parts[6]

	canBecomeRoot := false
	if userCanLogin(shell) {
		if utils.IsRoot() {
			if utils.CheckBinInstalled(logger, "sudo") {
				var err error
				canBecomeRoot, err = checkSudoRoot(username, logger)
				if err != nil {
					logger.Errorf("Unable to check if user %s can be root: %v", username, err)
				}
			} else {
				logger.Debugf("sudo not installed; skipping root check for user %s", username)
			}
		}
	}

	return &models.User{
		Username:      username,
		Uid:           uid,
		Gid:           gid,
		HomeDir:       home,
		Shell:         shell,
		CanBecomeRoot: canBecomeRoot,
		Name:          parts[4],
	}, nil
}

// userCanLogin check if user can login
func userCanLogin(shell string) bool {
	nologinShells := []string{
		"/sbin/nologin",
		"/usr/sbin/nologin",
		"/bin/false",
		"",
	}

	return !slices.Contains(nologinShells, shell)
}

// checkSudoRoot check if user can be root in running sudo -l -U command
func checkSudoRoot(user string, logger *logrus.Logger) (bool, error) {
	checkBecomeRoot, err := exec.Command("sudo", "-l", "-U", user).CombinedOutput()
	if err != nil {
		logger.Errorf("Error during execution of sudo -l -U %s: %s, %s", user, checkBecomeRoot, err)
		return false, err
	}

	userCanBecomeRoot := strings.Contains(string(checkBecomeRoot), "(ALL : ALL) ALL") ||
		strings.Contains(string(checkBecomeRoot), "(ALL) ALL")

	return userCanBecomeRoot, nil
}
