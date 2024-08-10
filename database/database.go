package database

import (
	"fmt"
	"time"
)

type Database struct {
	data map[string]interface{}
}

func New() *Database {
	return &Database{
		data: make(map[string]interface{}),
	}
}

func (db *Database) Value(key string) (interface{}, error) {
	// simulate 500ms roundtrip to the distributed cache
	time.Sleep(500 * time.Millisecond)

	return db.data[key], nil
}

func (db *Database) Store(key string, value interface{}) error {
	// simulate 500ms roundtrip to the distributed cache
	time.Sleep(500 * time.Millisecond)

	db.data[key] = value
	return nil
}

func FillDatabase(db *Database) {
	for i := 0; i < 10; i++ {
		db.Store(fmt.Sprintf("key%d", i),
			fmt.Sprintf("value%d", i))
	}
}
