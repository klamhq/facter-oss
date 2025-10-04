package users

import (
	"context"
	"testing"

	"github.com/klamhq/facter-oss/pkg/options"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	cfg := options.RunOptions{}
	res := New(logrus.New(), &cfg.Facter.Inventory.User)
	assert.NotNil(t, res)
}

func TestCollectUsersNoPasswdFile(t *testing.T) {
	cfg := options.RunOptions{}
	u := New(logrus.New(), &cfg.Facter.Inventory.User)
	ctx := context.Background()
	res, err := u.CollectUsers(ctx)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestCollectUsers(t *testing.T) {
	cfg := options.RunOptions{}
	cfg.Facter.Inventory.User.PasswdFile = "../../collectors/tests/passwd"
	u := New(logrus.New(), &cfg.Facter.Inventory.User)
	ctx := context.Background()
	res, err := u.CollectUsers(ctx)
	assert.Nil(t, err)
	assert.NotEmpty(t, res)
	assert.Len(t, res, 20)
	assert.Contains(t, res[19].Username, "bob")
}
