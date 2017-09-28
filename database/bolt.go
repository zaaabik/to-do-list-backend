package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"strconv"
	"time"
)

const dbBucket = "toDoList"

type BoltDb struct {
	db *bolt.DB
}

func NewBoltDb(path string) (*BoltDb, error) {
	if path == "" {
		path = "database.db"
	}
	db, err := bolt.Open(path, 0600, nil)
	return &BoltDb{db}, err
}

func (b *BoltDb) Save(item ToDo) (string, error) {
	key := strconv.FormatInt(time.Now().Unix(), 10)
	log.Print(key)
	enc, err := json.Marshal(item)
	if err != nil {
		return "", err
	}

	err = b.db.Update(func(tx *bolt.Tx) error {
		req, err := tx.CreateBucketIfNotExists([]byte(dbBucket))

		if err != nil {
			return err
		}

		err = req.Put([]byte(key), enc)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return key, nil
}
func (b *BoltDb) GetAll() ([]ToDotWithId, error) {
	var res []ToDotWithId
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucket))
		if b == nil {
			var emptyResult []ToDotWithId
			res = emptyResult
			return nil
		}
		err := b.ForEach(func(k, v []byte) error {
			var item ToDotWithId
			json.Unmarshal(v, &item)
			item.Id = string(k)
			res = append(res, item)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	fmt.Print(res)
	return res, nil
}

func (b *BoltDb) DeleteAll() error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucket))
		if b == nil {
			return nil
		}
		tx.DeleteBucket([]byte(dbBucket))
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *BoltDb) Delete(key string) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucket))
		if b == nil {
			return errors.New("bucket is empty=" + key)
		}
		err := b.Delete([]byte(key))
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *BoltDb) Close() {
	b.db.Close()
}
