package packages

import (
	"context"

	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
)

type PackagesCollector interface {
	CollectPackages(ctx context.Context) ([]*schema.Package, error)
}
