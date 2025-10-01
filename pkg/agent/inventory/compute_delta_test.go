package inventory

import (
	"testing"
	"time"

	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestComputeDelta(t *testing.T) {
	oldInv := &schema.HostInventory{
		Hostname:  "test-host",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		Packages: []*schema.Package{
			{Name: "pkg1", Version: "1.0.0"},
			{Name: "pkg2", Version: "1.0.0"},
		},
		Users: []*schema.User{
			{Username: "user1"},
		},
	}

	newInv := &schema.HostInventory{
		Hostname:  "test-host",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		Packages: []*schema.Package{
			{Name: "pkg1", Version: "1.0.0"},
			{Name: "pkg3", Version: "1.0.0"},
		},
		Users: []*schema.User{
			{Username: "user1"},
			{Username: "user2"},
		},
	}

	logger := logrus.New()
	delta := ComputeDelta(oldInv, newInv, logger)

	assert.NotNil(t, delta)
	assert.Equal(t, "test-host", delta.Hostname)
	assert.NotEmpty(t, delta.UpdatedAt)
	assert.Len(t, delta.PackagesAdded, 1)
	assert.Len(t, delta.PackagesRemoved, 1)
	assert.Len(t, delta.UsersAdded, 1)
	assert.Len(t, delta.UsersRemoved, 0)
}

func TestIsDeltaEmpty(t *testing.T) {
	delta := &schema.HostDeltaInventory{
		PackagesAdded:   []*schema.Package{},
		PackagesRemoved: []*schema.Package{},
		UsersAdded:      []*schema.User{},
		UsersRemoved:    []*schema.User{},
	}

	assert.True(t, IsDeltaEmpty(delta))

	delta.PackagesAdded = append(delta.PackagesAdded, &schema.Package{Name: "pkg1"})
	assert.False(t, IsDeltaEmpty(delta))
}

func TestDiffByHash(t *testing.T) {
	oldList := []*schema.Package{
		{Name: "pkg1", Version: "1.0.0"},
		{Name: "pkg2", Version: "1.0.0"},
	}

	newList := []*schema.Package{
		{Name: "pkg1", Version: "1.0.0"},
		{Name: "pkg3", Version: "1.0.0"},
	}

	added, removed, _ := DiffGenericByHash(oldList, newList, func(p *schema.Package) string { return p.Name })

	assert.Len(t, added, 1)
	assert.Equal(t, "pkg3", added[0].Name)
	assert.Len(t, removed, 1)
	assert.Equal(t, "pkg2", removed[0].Name)
}
