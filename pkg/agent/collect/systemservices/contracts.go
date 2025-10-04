package systemservices

import (
	"context"

	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
)

type SystemServicesCollector interface {
	CollectSystemServices(ctx context.Context, initsystem string) ([]*schema.SystemdService, error)
}
