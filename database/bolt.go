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
	id := time.Now().UnixNano() / int64(time.Millisecond)
	key := strconv.FormatInt(id, 10)
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
func (b *BoltDb) Get(id string) (ToDo, error) {
	var res ToDo
	emptyResult := ToDo{}
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucket))
		if b == nil {
			res = emptyResult
			return nil
		}
		rawToDo := b.Get([]byte(id))
		err := json.Unmarshal(rawToDo, &res)
		if err != nil {
			return err
		}
		res.Id = id
		return nil
	})

	if err != nil {
		return emptyResult, err
	}
	fmt.Print(res)
	return res, nil
}

func (b *BoltDb) GetAll() ([]ToDo, error) {
	var res []ToDo
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucket))
		if b == nil {
			var emptyResult []ToDo
			res = emptyResult
			return nil
		}
		err := b.ForEach(func(k, v []byte) error {
			var item ToDo
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

func (b *BoltDb) CloseTodo(id string) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucket))
		if b == nil {
			return errors.New("bucket is empty=" + id)
		}
		v := b.Get([]byte(id))
		var item ToDo
		json.Unmarshal(v, &item)
		item.Id = string(id)
		item.IsClosed = true
		item.UpdatedAt = time.Now().Format(time.ANSIC)
		data, err := json.Marshal(item)
		err = b.Put([]byte(id), data)
		return err
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
