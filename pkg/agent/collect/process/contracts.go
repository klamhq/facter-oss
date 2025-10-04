package process

import (
	"context"

	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
)

type ProcessCollector interface {
	CollectProcess(ctx context.Context) ([]*schema.Process, error)
}
