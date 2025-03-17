package dbfile

import (
	"encoding/json"
	"os"
	"sync"
)

// KeyValueDBFile - in file database
type KeyValueDBFile interface {
	Put(key string, val []byte) error
	Remove(key string) error
	Keys() []string
	Get(key string) []byte
}

type keyFile struct {
	data map[string][]byte
	file string
	lock sync.Mutex
}

// Put - put data into database under key
func (k *keyFile) Put(key string, val []byte) error {
	k.lock.Lock()
	defer k.lock.Unlock()
	k.data[key] = val
	return k.save()
}

// Remove - remove key
func (k *keyFile) Remove(key string) error {
	k.lock.Lock()
	defer k.lock.Unlock()
	delete(k.data, key)
	return k.save()
}

// Keys - list of keys in database
func (k *keyFile) Keys() []string {
	k.lock.Lock()
	defer k.lock.Unlock()
	keys := make([]string, len(k.data))
	i := 0
	for ek, _ := range k.data {
		keys[i] = ek
		i++
	}
	return keys
}

// Get - get data from database
func (k *keyFile) Get(key string) []byte {
	buff := k.data[key]
	return buff
}

func (k *keyFile) save() error {
	buff, err := json.Marshal(k.data)
	if err != nil {
		return err
	}
	return os.WriteFile(k.file, buff, os.ModePerm)
}

func (k *keyFile) load() error {
	buff, err := os.ReadFile(k.file)
	if err != nil {
		return err
	}
	m := make(map[string][]byte)
	err = json.Unmarshal(buff, &m)
	if err != nil {
		return err
	}
	k.data = m
	return nil
}

// NewKeyValueDBFile - create new in file database
func NewKeyValueDBFile(file string) (KeyValueDBFile, error) {
	kf := &keyFile{
		data: make(map[string][]byte),
		file: file,
		lock: sync.Mutex{},
	}
	if _, err := os.Stat(file); err == nil {
		kf.lock.Lock()
		defer kf.lock.Unlock()
		err := kf.load()
		if err != nil {
			return nil, err
		}
	}
	return kf, nil
}
