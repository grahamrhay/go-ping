package main

import (
	"bytes"
	"encoding/gob"
	"github.com/steveyen/gkvlite"
	"log"
	"os"
)

type Store struct {
	f *os.File
	s *gkvlite.Store
}

func openStore() (*Store, error) {
	log.Println("Opening store")
	f, err := os.OpenFile("./db", os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return nil, err
	}
	s, err := gkvlite.NewStore(f)
	if err != nil {
		return nil, err
	}
	return &Store{f: f, s: s}, nil
}

func closeStore(store *Store) {
	log.Println("Closing store")
	store.s.Close()
	store.f.Close()
}

func writeToStore(store *Store, coll string, item interface{}, key string) error {
	log.Printf("Writing item to store. Coll: %v, Key: %v, Item: %v\n", coll, key, item)
	c := store.s.GetCollection(coll)
	if c == nil {
		log.Println("Collection doesn't exist, creating it")
		c = store.s.SetCollection(coll, nil)
	}
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(item)
	if err != nil {
		return err
	}
	err = c.Set([]byte(key), buffer.Bytes())
	if err != nil {
		return err
	}
	err = store.s.Flush()
	if err != nil {
		return err
	}
	return store.f.Sync()
}

func getFromStore(store *Store, coll string, key string) (*PingResult, error) {
	log.Printf("Retrieving item from store. Coll: %v, Key: %v\n", coll, key)
	c := store.s.GetCollection(coll)
	if c == nil {
		log.Println("Collection doesn't exist")
		return nil, nil
	}
	itemBytes, err := c.Get([]byte(key))
	if err != nil {
		return nil, err
	}
	if itemBytes == nil {
		return nil, nil
	}
	buffer := bytes.NewBuffer(itemBytes)
	dec := gob.NewDecoder(buffer)
	var result PingResult
	err = dec.Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
