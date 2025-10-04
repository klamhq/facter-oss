package users

import (
	"context"

	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
)

type UsersCollector interface {
	CollectUsers(ctx context.Context) ([]*schema.User, error)
}
