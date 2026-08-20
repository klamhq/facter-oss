package store

import (
	"path"
	"testing"
	"time"

	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/stretchr/testify/assert"
)

func TestNewBoltInventoryStore(t *testing.T) {
	path := path.Join(t.TempDir(), "test.db")
	store, err := NewBoltInventoryStore(path)
	assert.NotNil(t, store)
	assert.NoError(t, err)
}

func TestSaveAndGetInventory(t *testing.T) {
	path := path.Join(t.TempDir(), "test.db")
	store, err := NewBoltInventoryStore(path)
	assert.NoError(t, err)

	inv := &schema.HostInventory{
		Hostname:  "test-host",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		Packages:  []*schema.Package{{Name: "pkg1", Version: "1.0.0"}},
		Users:     []*schema.User{{Username: "user1"}},
	}

	err = store.Save("test-host", inv)
	assert.NoError(t, err)

	retrieved, err := store.Get("test-host")
	assert.NoError(t, err)

	assert.Equal(t, inv.Users[0].Username, retrieved.Users[0].Username, "Usernames should match")
	assert.Equal(t, inv.Users[0].Uid, retrieved.Users[0].Uid, "Uids should match")
	assert.Equal(t, inv.Users[0].Gid, retrieved.Users[0].Gid, "Gids should match")
	assert.Equal(t, inv.Users[0].Home, retrieved.Users[0].Home, "Homes should match")
	assert.Equal(t, inv.Users[0].Shell, retrieved.Users[0].Shell, "Shells should match")
	assert.Equal(t, inv.Users[0].CanBecomeRoot, retrieved.Users[0].CanBecomeRoot, "Canbecomeroot should match")
	assert.Equal(t, inv.Users[0].UpdatedAt, retrieved.Users[0].UpdatedAt, "UpdatedAt should match")

	assert.Equal(t, inv.Hostname, retrieved.Hostname, "Hostnames should match")
	assert.Equal(t, inv.CreatedAt, retrieved.CreatedAt, "createdAt should match")
	assert.Len(t, retrieved.Packages, 1, "Should have one package")
	assert.Equal(t, inv.Packages[0].Name, retrieved.Packages[0].Name, "Package names should match")
	assert.Equal(t, inv.Packages[0].Version, retrieved.Packages[0].Version, "Package versions should match")
}

func TestDeleteInventory(t *testing.T) {
	path := path.Join(t.TempDir(), "test.db")
	store, err := NewBoltInventoryStore(path)
	assert.NoError(t, err)

	inv := &schema.HostInventory{
		Hostname:  "test-host",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		Packages:  []*schema.Package{{Name: "pkg1", Version: "1.0.0"}},
		Users:     []*schema.User{{Username: "user1"}},
	}

	err = store.Save("test-host", inv)
	assert.NoError(t, err)

	err = store.Delete("test-host")
	assert.Nil(t, err, "Retrieved inventory should be nil after deletion")
	assert.NoError(t, err)
	retrieved, err := store.Get("test-host")
	assert.Error(t, err, "Expected error when retrieving deleted inventory")
	assert.EqualError(t, err, "not found", "Expected 'not found' error")
	assert.NotNil(t, retrieved, "Retrieved inventory should not be nil after deletion")

}

func TestCloseStore(t *testing.T) {
	path := path.Join(t.TempDir(), "test.db")
	store, err := NewBoltInventoryStore(path)
	assert.NoError(t, err)
	err = store.Close()
	assert.NoError(t, err)
}

func TestSaveRevision_IdempotentByRevisionID(t *testing.T) {
	path := path.Join(t.TempDir(), "revisions.db")
	store, err := NewBoltInventoryStore(path)
	assert.NoError(t, err)

	rev := &schema.InventoryRevisionEnvelope{
		RevisionId:         "rev-1",
		Sequence:           1,
		PreviousRevisionId: "",
		PreviousSequence:   0,
		Hostname:           "test-host",
		SourceType:         schema.SourceType_SOURCE_TYPE_FULL,
		StateHash:          "abc123",
		Payload:            &schema.InventoryRevisionEnvelope_Full{Full: &schema.HostInventory{Hostname: "test-host"}},
	}

	assert.NoError(t, store.SaveRevision("test-host", rev))
	assert.NoError(t, store.SaveRevision("test-host", rev))

	got, err := store.GetRevision("test-host")
	assert.NoError(t, err)
	assert.Equal(t, "rev-1", got.RevisionId)
	assert.EqualValues(t, 1, got.Sequence)
}

func TestCloseNilStore(t *testing.T) {
	var store *boltInventoryStore
	err := store.Close()
	assert.NoError(t, err)
}
