// pkg/collect/platform/contracts.go
package platform

import (
	"context"

	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
)

type PlatformCollector interface {
	CollectPlatform(ctx context.Context) (*schema.Platform, error)
}
