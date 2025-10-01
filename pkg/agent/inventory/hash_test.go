package inventory

import (
	"testing"

	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	assert := assert.New(t)
	var pkg schema.Package
	pkg.Architecture = "amd64"
	pkg.Name = "test-package"
	pkg.Version = "1.0.0"

	h := hashMessage(&pkg)
	assert.NotNil(h)
	assert.NotEqual(uint64(0), h)
	assert.Equal(uint64(2399281260310825189), h) // Replace with expected hash value
}

func TestResetUpdatedAtField(t *testing.T) {
	assert := assert.New(t)
	var inv schema.HostDeltaInventory
	inv.UpdatedAt = "2023-10-01T12:00:00Z"

	resetUpdatedAtField(&inv)

	assert.Equal("", inv.UpdatedAt, "UpdatedAt field should be reset to empty string")
}
