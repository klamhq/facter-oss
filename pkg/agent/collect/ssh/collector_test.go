package ssh

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/klamhq/facter-oss/pkg/options"
	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const mockKeyRSA = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCQh08WbnI8YIATA/frSJGvW2oUjvZ387QDwH5OPfw8gN+SGpcWYdXBMVljVh72zPjiEBkeU9MXcEccKMGQXCE6AG5xsxA9saXJ3JIB/Ydo1wgabWXw3/HbwE41PWNu0l4xcM3d2IJRMxqk1HB5ipwvWkbZ9Xs3XqngYIblHKkdkCAq4cZ6EAWA2LgvNcUZh31IN0m3d6ZcOF4Xgg/qCIEoiPXM2GDRttASXjLNpEiB1sF/lbKCEyzmrejNUMEOiw/Hlwid+y5Vmmefn7sKVoCZ42ZdDuynzFWR+fs4ISjMoJvPKrtFoea+JMdO9hNp4QUuwbYDh8CJNTq1pV/8UTJOFk4FShrncYnWlgeRmKp+P39QE/JFJMwxErkGT/mTH38l9QYIgeflNpsOTcM3fpprSxU/sVIxb3IFssQQwNQQ9Bp+1eo2nwrjxuY1EyZiqfUsAo+OdRpcQwV+2LyNneNlp7Az3n5/xuYgwTwKUslPdhlA4h/mjfdF9EAIl20CwZM= user@host`

func TestNew(t *testing.T) {
	cfg := options.RunOptions{}
	res := New(logrus.New(), &cfg.Facter.Inventory.SSH)
	assert.NotNil(t, res)
}

func TestCollectSSH(t *testing.T) {
	cfg := options.RunOptions{}
	ctx := context.Background()
	dir := t.TempDir()
	dirSsh := dir + "/.ssh"
	os.Mkdir(dirSsh, 0744)
	// Create a known_hosts file with a valid entry
	knownHostsFile := filepath.Join(dirSsh, "known_hosts")
	knownHostsContent := "host.example.com " + mockKeyRSA
	err := os.WriteFile(knownHostsFile, []byte(knownHostsContent), 0644)
	assert.NoError(t, err)

	pubFile := filepath.Join(dirSsh, "id_rsa.pub")
	content := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCQh08WbnI8YIATA/frSJGvW2oUjvZ387QDwH5OPfw8gN+SGpcWYdXBMVljVh72zPjiEBkeU9MXcEccKMGQXCE6AG5xsxA9saXJ3JIB/Ydo1wgabWXw3/HbwE41PWNu0l4xcM3d2IJRMxqk1HB5ipwvWkbZ9Xs3XqngYIblHKkdkCAq4cZ6EAWA2LgvNcUZh31IN0m3d6ZcOF4Xgg/qCIEoiPXM2GDRttASXjLNpEiB1sF/lbKCEyzmrejNUMEOiw/Hlwid+y5Vmmefn7sKVoCZ42ZdDuynzFWR+fs4ISjMoJvPKrtFoea+JMdO9hNp4QUuwbYDh8CJNTq1pV/8UTJOFk4FShrncYnWlgeRmKp+P39QE/JFJMwxErkGT/mTH38l9QYIgeflNpsOTcM3fpprSxU/sVIxb3IFssQQwNQQ9Bp+1eo2nwrjxuY1EyZiqfUsAo+OdRpcQwV+2LyNneNlp7Az3n5/xuYgwTwKUslPdhlA4h/mjfdF9EAIl20CwZM= user@host"
	err = os.WriteFile(pubFile, []byte(content), 0644)
	assert.NoError(t, err)

	authFile := filepath.Join(dirSsh, "authorized_keys")
	content1 := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCQh08WbnI8YIATA/frSJGvW2oUjvZ387QDwH5OPfw8gN+SGpcWYdXBMVljVh72zPjiEBkeU9MXcEccKMGQXCE6AG5xsxA9saXJ3JIB/Ydo1wgabWXw3/HbwE41PWNu0l4xcM3d2IJRMxqk1HB5ipwvWkbZ9Xs3XqngYIblHKkdkCAq4cZ6EAWA2LgvNcUZh31IN0m3d6ZcOF4Xgg/qCIEoiPXM2GDRttASXjLNpEiB1sF/lbKCEyzmrejNUMEOiw/Hlwid+y5Vmmefn7sKVoCZ42ZdDuynzFWR+fs4ISjMoJvPKrtFoea+JMdO9hNp4QUuwbYDh8CJNTq1pV/8UTJOFk4FShrncYnWlgeRmKp+P39QE/JFJMwxErkGT/mTH38l9QYIgeflNpsOTcM3fpprSxU/sVIxb3IFssQQwNQQ9Bp+1eo2nwrjxuY1EyZiqfUsAo+OdRpcQwV+2LyNneNlp7Az3n5/xuYgwTwKUslPdhlA4h/mjfdF9EAIl20CwZM= user@host"
	err = os.WriteFile(authFile, []byte(content1), 0644)
	assert.NoError(t, err)
	users := []*schema.User{{
		Username: "a",
		Home:     dir,
	}}

	s := New(logrus.New(), &cfg.Facter.Inventory.SSH)
	ska, kh, ski, err := s.CollectSSHInfos(ctx, users)
	assert.NoError(t, err)
	assert.NotNil(t, ska)
	assert.Equal(t, ska[0].Fingerprint, "e3:3e:53:ec:94:22:41:f6:81:cc:75:ab:ed:ee:50:19")
	assert.NotNil(t, kh)
	assert.Equal(t, kh[0].Hostname, "host.example.com")
	assert.Equal(t, kh[0].Type, "ssh-rsa")
	assert.Equal(t, kh[0].Fingerprint, "SHA256:onumzaYU98XSgnDrr1a2yBjuArXo+EbG/MMRA/bhtG8")

	assert.NotNil(t, ski)

	assert.Equal(t, ski[0].Fingerprint, "SHA256:onumzaYU98XSgnDrr1a2yBjuArXo+EbG/MMRA/bhtG8")
	assert.Equal(t, ski[0].Type, "ssh-rsa")
	assert.Equal(t, ski[0].Length, int64(3072))
	assert.Equal(t, ski[0].Path, dir)
	assert.Equal(t, ski[0].Name, "authorized_keys")
	assert.Equal(t, ski[0].Comment, "user@host")
	assert.Equal(t, ski[0].FromAuthorizedKeys, true)

}
