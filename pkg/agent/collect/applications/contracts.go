package applications

import (
	"context"

	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
)

type ApplicationsCollector interface {
	Name() string
	CollectApplications(ctx context.Context) ([]*schema.Application, error)
}
