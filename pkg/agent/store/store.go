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
	SaveRevision(hostname string, rev *schema.InventoryRevisionEnvelope) error
	GetRevision(hostname string) (*schema.InventoryRevisionEnvelope, error)
	DeleteRevision(hostname string) error
	Close() error
}

type boltInventoryStore struct {
	db *bolt.DB
}

const (
	inventoryBucket = "inventory"
	revisionBucket  = "revisions"
)

func NewBoltInventoryStore(path string) (*boltInventoryStore, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(inventoryBucket)); err != nil {
			return err
		}
		_, err := tx.CreateBucketIfNotExists([]byte(revisionBucket))
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

func (b *boltInventoryStore) SaveRevision(hostname string, rev *schema.InventoryRevisionEnvelope) error {
	if rev == nil {
		return nil
	}
	current, err := b.GetRevision(hostname)
	if err == nil && current != nil && current.RevisionId == rev.RevisionId {
		return nil
	}
	data, err := proto.Marshal(rev)
	if err != nil {
		return err
	}
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(revisionBucket))
		return bucket.Put([]byte(hostname), data)
	})
}

func (b *boltInventoryStore) GetRevision(hostname string) (*schema.InventoryRevisionEnvelope, error) {
	var rev schema.InventoryRevisionEnvelope
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(revisionBucket))
		if bucket == nil {
			return fmt.Errorf("revision bucket not found")
		}
		data := bucket.Get([]byte(hostname))
		if data == nil {
			return fmt.Errorf("data not found")
		}
		return proto.Unmarshal(data, &rev)
	})
	if err != nil {
		return nil, err
	}
	return &rev, err
}

func (b *boltInventoryStore) DeleteRevision(hostname string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(revisionBucket))
		return bucket.Delete([]byte(hostname))
	})
}

func (b *boltInventoryStore) Close() error {
	if b == nil || b.db == nil {
		return nil
	}
	return b.db.Close()
}
