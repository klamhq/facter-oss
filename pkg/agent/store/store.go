package store

import (
	"fmt"

	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	bolt "go.etcd.io/bbolt"
	proto "google.golang.org/protobuf/proto"
)

type InventoryStore interface {
	Get(hostname string) (*schema.HostInventory, error)
	Save(hostname string, inv *schema.HostInventory) error
	Delete(hostname string) error
	Close() error
}

type boltInventoryStore struct {
	db *bolt.DB
}

const inventoryBucket = "inventory"

func NewBoltInventoryStore(path string) (*boltInventoryStore, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	// init bucket
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(inventoryBucket))
		return err
	})
	return &boltInventoryStore{db}, err
}

func (b *boltInventoryStore) Save(hostname string, inv *schema.HostInventory) error {
	data, err := proto.Marshal(inv)
	if err != nil {
		return err
	}
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(inventoryBucket))
		return bucket.Put([]byte(hostname), data)
	})
}

func (b *boltInventoryStore) Get(hostname string) (*schema.HostInventory, error) {
	var inv schema.HostInventory
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(inventoryBucket))
		data := bucket.Get([]byte(hostname))
		if data == nil {
			return fmt.Errorf("not found")
		}
		return proto.Unmarshal(data, &inv)
	})
	return &inv, err
}

func (b *boltInventoryStore) Delete(hostname string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(inventoryBucket))
		return bucket.Delete([]byte(hostname))
	})
}

func (b *boltInventoryStore) Close() error {
	if b == nil || b.db == nil {
		return nil
	}
	return b.db.Close()
}
