package networks

import (
	"context"

	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
)

type NetworksCollector interface {
	CollectNetworks(ctx context.Context) (*schema.Network, error)
}
