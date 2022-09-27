package db

import (
	"fmt"
	"os"

	bolt "go.etcd.io/bbolt"
)

// This is a wrapper around the bolt client
type Bolt struct {
	Db *bolt.DB
}

func (b *Bolt) Close() {
	b.Db.Close()
}

// Takes a relative path to a bolt db file
// - opens the db
// - updates the Bolt.Db field with a reference to the bolt instance
// - returns the bolt instance as a courtesy
func (b *Bolt) Init(dbPath string) *bolt.DB {
	if dbPath == "" {
		fmt.Println("Please provide a path to a bolt database")
		os.Exit(1)
	}

	db, err := bolt.Open(dbPath, 0666, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Note, do not call close here.
	// defer db.Close()

	b.Db = db
	return db
}

func (b *Bolt) ListBuckets() []string {
	bucketSlice := []string{}

	b.Db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			bucketSlice = append(bucketSlice, string(name))
			return nil
		})
	})

	return bucketSlice
}

func (b *Bolt) ListKV(bucketName string) []map[string]interface{} {
	kvSlice := []map[string]interface{}{}

	b.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		return bucket.ForEach(func(k, v []byte) error {
			kvSlice = append(kvSlice, map[string]interface{}{string(k): string(v)})
			return nil
		})
	})

	return kvSlice
}

func (b *Bolt) DescribeBucket(name string) bolt.BucketStats {
	var stats bolt.BucketStats

	b.Db.View(func(tx *bolt.Tx) error {
		stats = tx.Bucket([]byte(name)).Stats()
		return nil
	})

	return stats
}
