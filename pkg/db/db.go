package db

import (
	"log"

	"go.etcd.io/bbolt"
)

var db *bbolt.DB

func InitDB(dbPath string) {
	var err error
	db, err = bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatalf("Could not open db: %v", err)
	}
}

func CloseDB() {
	db.Close()
}
