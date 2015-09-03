// Copyright 2015 The Loadcat Authors. All rights reserved.

package data

import (
	"github.com/boltdb/bolt"
)

var DB *bolt.DB

func OpenDB(path string) error {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return err
	}
	DB = db
	return nil
}

func InitDB() error {
	return DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("balancers"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("servers"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("settings"))
		if err != nil {
			return err
		}
		return nil
	})
}
