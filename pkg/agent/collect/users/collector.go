package users

import (
	"context"

	"github.com/klamhq/facter-oss/pkg/agent/collectors/users"
	"github.com/klamhq/facter-oss/pkg/options"
	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/sirupsen/logrus"
)

type UserCollectorImpl struct {
	log *logrus.Logger
	cfg *options.UserOptions
}

func New(log *logrus.Logger, cfg *options.UserOptions) *UserCollectorImpl {

	return &UserCollectorImpl{
		log: log,
		cfg: cfg,
	}
}

func (c *UserCollectorImpl) Name() string { return "users" }

func (c *UserCollectorImpl) CollectUsers(ctx context.Context) ([]*schema.User, error) {
	c.log.Info("Crafting users")
	var collectedUsers []*schema.User
	getConnectedUsers := users.GetConnectedUsers(c.log)

	getUsers, err := users.GetSystemUsers(c.cfg.PasswdFile, c.log)

	Mergeusers := users.MergeUsersAndSessions(getUsers, getConnectedUsers)

	if err != nil {
		return nil, err
	}

	for _, u := range Mergeusers {
		var sessions []*schema.Session
		for _, s := range u.Session {
			sessions = append(sessions, &schema.Session{
				Terminal:  s.Terminal,
				Host:      s.Host,
				Started:   s.Started,
				Connected: s.Connected,
			})
		}
		protoUser := schema.User{
			Username:      u.Username,
			Uid:           u.Uid,
			Gid:           u.Gid,
			Home:          u.HomeDir,
			Sessions:      sessions,
			Shell:         u.Shell,
			CanBecomeRoot: u.CanBecomeRoot,
		}

		collectedUsers = append(collectedUsers, &protoUser)
	}

	return collectedUsers, nil
}
