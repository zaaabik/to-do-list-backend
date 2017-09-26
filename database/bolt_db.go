package database

import (
	_ "encoding/binary"
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

func (b *BoltDb) Save(item ListItem) (string, error) {
	t := strconv.FormatInt(time.Now().Unix(), 10)
	log.Print(t)
	enc, err := json.Marshal(item)
	if err != nil {
		return "", err
	}

	err = b.db.Update(func(tx *bolt.Tx) error {
		req, err := tx.CreateBucketIfNotExists([]byte(dbBucket))

		if err != nil {
			return err
		}

		err = req.Put([]byte(t), enc)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return t, nil
}
func (b *BoltDb) GetAll() ([]ListItemWithId, error) {
	var res []ListItemWithId
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucket))
		if b == nil {
			var emptyResult []ListItemWithId
			res = emptyResult
			return nil
		}
		err := b.ForEach(func(k, v []byte) error {
			var item ListItemWithId
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
