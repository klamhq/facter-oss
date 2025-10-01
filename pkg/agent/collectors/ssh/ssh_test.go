package ssh_test

import (
	"testing"

	"github.com/klamhq/facter-oss/pkg/agent/collectors/ssh"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// Mock data
const mockKeyRSA = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCQh08WbnI8YIATA/frSJGvW2oUjvZ387QDwH5OPfw8gN+SGpcWYdXBMVljVh72zPjiEBkeU9MXcEccKMGQXCE6AG5xsxA9saXJ3JIB/Ydo1wgabWXw3/HbwE41PWNu0l4xcM3d2IJRMxqk1HB5ipwvWkbZ9Xs3XqngYIblHKkdkCAq4cZ6EAWA2LgvNcUZh31IN0m3d6ZcOF4Xgg/qCIEoiPXM2GDRttASXjLNpEiB1sF/lbKCEyzmrejNUMEOiw/Hlwid+y5Vmmefn7sKVoCZ42ZdDuynzFWR+fs4ISjMoJvPKrtFoea+JMdO9hNp4QUuwbYDh8CJNTq1pV/8UTJOFk4FShrncYnWlgeRmKp+P39QE/JFJMwxErkGT/mTH38l9QYIgeflNpsOTcM3fpprSxU/sVIxb3IFssQQwNQQ9Bp+1eo2nwrjxuY1EyZiqfUsAo+OdRpcQwV+2LyNneNlp7Az3n5/xuYgwTwKUslPdhlA4h/mjfdF9EAIl20CwZM= user@host`

func TestGetSshKeyInfo_ValidKey(t *testing.T) {
	keyList := []string{mockKeyRSA}
	keyPaths := []string{"/home/user/.ssh/id_rsa.pub"}
	logger := &logrus.Logger{}
	result := ssh.GetSshKeyInfo(logger, keyList, keyPaths)

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
	result := ssh.GetSshKeyInfo(logger, keyList, keyPaths)

	assert.Len(t, result, 0)
}

func TestGetPath(t *testing.T) {
	base, file := ssh.GetPath("/home/user/.ssh/id_rsa.pub")

	assert.Equal(t, "/home/user", base)
	assert.Equal(t, "id_rsa.pub", file)
}

func TestContains(t *testing.T) {
	values := []string{"one", "two", "three"}

	assert.True(t, ssh.Contains(values, "two"))
	assert.False(t, ssh.Contains(values, "four"))
}
