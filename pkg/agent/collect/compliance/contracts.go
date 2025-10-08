package compliance

import (
	"context"

	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
)

type ComplianceCollector interface {
	CollectCompliance(ctx context.Context) (*schema.ComplianceReport, error)
}
