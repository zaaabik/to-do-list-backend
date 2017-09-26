package database

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"time"
	"fmt"
	"log"
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

func (b *BoltDb) Save(item ListItem) error {
	enc, err := json.Marshal(item)
	if err != nil {
		return err
	}

	err = b.db.Update(func(tx *bolt.Tx) error {
		req, err := tx.CreateBucketIfNotExists([]byte(dbBucket))

		if err != nil {
			return err
		}

		err = req.Put([]byte(time.Now().Format(time.UnixDate)), enc)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func (b *BoltDb) GetAll() ([]ListItem,error){
	var res []ListItem
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucket))
		if b == nil {
			var emptyResult []ListItem
			res = emptyResult
			return nil
		}
		err := b.ForEach(func(k, v []byte) error {
			var item ListItem
			json.Unmarshal(v,&item)
			res = append(res, item)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil,err
	}
	fmt.Print(res)
	return res, nil
}

func (b *BoltDb) DeleteAll() error{
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

func (b *BoltDb) Delete(item ListItem) error {
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

func (b *BoltDb) Close() {
	b.db.Close()
}
