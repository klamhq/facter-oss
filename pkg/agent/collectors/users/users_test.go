package users

import (
	"os"
	"testing"

	"github.com/klamhq/facter-oss/pkg/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const passwdFilename = "../tests/passwd"

func TestGetSystemUsers(t *testing.T) {
	logger := logrus.New()
	users, err := GetSystemUsers(passwdFilename, logger)

	assert.NoError(t, err)
	assert.NotEmpty(t, users)
	// Ajuste ce nombre selon ton fichier de test réel
	assert.True(t, len(users) == 20, "expected 20 users")
}

func TestGetSystemUsers_Failure(t *testing.T) {
	logger := logrus.New()
	users, err := GetSystemUsers("/invalid/path", logger)

	assert.Error(t, err)
	assert.Empty(t, users)
}

func TestParseUserLine_ValidLines(t *testing.T) {
	if _, isCI := os.LookupEnv("CI"); isCI {
		t.Skip("Skipping test in CI (sudo call may prompt)")
	}

	logger := logrus.New()

	tests := []struct {
		name string
		line string
		want models.User
	}{
		{
			name: "simple root user",
			line: "root:x:0:0:root:/root:/bin/bash",
			want: models.User{
				Username: "root",
				Uid:      "0",
				Gid:      "0",
				HomeDir:  "/root",
				Shell:    "/bin/bash",
			},
		},
		{
			name: "another root format",
			line: "root:x:0:0::/root:/bin/zsh",
			want: models.User{
				Username: "root",
				Uid:      "0",
				Gid:      "0",
				HomeDir:  "/root",
				Shell:    "/bin/zsh",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			user, err := parseUserLine(tc.line, logger)
			assert.NoError(t, err)
			assert.NotNil(t, user)
			assert.Equal(t, tc.want.Username, user.Username)
			assert.Equal(t, tc.want.Uid, user.Uid)
			assert.Equal(t, tc.want.Gid, user.Gid)
			assert.Equal(t, tc.want.HomeDir, user.HomeDir)
			assert.Equal(t, tc.want.Shell, user.Shell)
		})
	}
}

func TestParseUserLine_InvalidLines(t *testing.T) {
	logger := logrus.New()

	t.Run("too few fields", func(t *testing.T) {
		line := "root:x:0"
		user, err := parseUserLine(line, logger)
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("empty username", func(t *testing.T) {
		line := ":x:0:0::/root:/bin/bash"
		user, err := parseUserLine(line, logger)
		assert.NoError(t, err)
		assert.Nil(t, user)
	})

	t.Run("malformed line (no error expected, only partial data)", func(t *testing.T) {
		line := "ro`ot:x:0:@:::"
		user, err := parseUserLine(line, logger)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "ro`ot", user.Username)
		assert.Equal(t, "0", user.Uid)
		assert.Equal(t, "@", user.Gid)
	})
}

func TestGetConnectedUsers(t *testing.T) {
	logger := logrus.New()
	connectedUser := GetConnectedUsers(logger)
	assert.NotNil(t, connectedUser)
}

func TestCheckSudoRootFail(t *testing.T) {
	logger := logrus.New()
	b, err := checkSudoRoot("user", logger)
	assert.False(t, b)
	assert.Error(t, err)
}

func TestMergeUsersAndSessions(t *testing.T) {
	logger := logrus.New()
	users, err := GetSystemUsers("/etc/passwd", logger)
	assert.NoError(t, err)
	connectedUser := GetConnectedUsers(logger)
	assert.NotNil(t, connectedUser)
	mergeUsers := MergeUsersAndSessions(users, connectedUser)
	assert.NotEmpty(t, mergeUsers)

}
