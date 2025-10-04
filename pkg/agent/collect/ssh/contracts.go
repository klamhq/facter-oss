package ssh

import (
	"context"

	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
)

type SSHInfosCollector interface {
	CollectSSHInfos(ctx context.Context, users []*schema.User) ([]*schema.SshKeyAccess, []*schema.KnownHost, []*schema.SshKeyInfo, error)
}
