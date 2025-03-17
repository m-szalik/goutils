package dbfile

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"slices"
	"testing"
)

func withKeyValueDBFile(t *testing.T, f func(t *testing.T, db KeyValueDBFile)) {
	fn := fmt.Sprintf("test-%d.tmp", rand.Intn(1000))
	db, err := NewKeyValueDBFile(fn)
	assert.NoError(t, err)
	f(t, db)
	err = os.Remove(fn)
	t.Cleanup(func() {
		_ = os.Remove(fn)
	})
}

func TestKeys(t *testing.T) {
	keys := []string{"k0", "k1", "k2", "ax"}
	withKeyValueDBFile(t, func(t *testing.T, db KeyValueDBFile) {
		for i, key := range keys {
			err := db.Put(key, []byte(fmt.Sprintf("value-%d", i)))
			assert.NoError(t, err)
		}
		dbKeys := db.Keys()
		slices.Sort(dbKeys)
		slices.Sort(keys)
		assert.Equal(t, keys, dbKeys)
	})
}

func TestRemoveKey(t *testing.T) {
	withKeyValueDBFile(t, func(t *testing.T, db KeyValueDBFile) {
		err := db.Put("key", []byte("xyz"))
		assert.NoError(t, err)
		dbKeys := db.Keys()
		assert.Equal(t, []string{"key"}, dbKeys)
		value := db.Get("key")
		assert.Equal(t, "xyz", string(value))
		err = db.Remove("key")
		assert.NoError(t, err)
		value = db.Get("key")
		assert.Nil(t, value)
		dbKeys = db.Keys()
		assert.Equal(t, []string{}, dbKeys)
	})
}

func TestOverrideKey(t *testing.T) {
	withKeyValueDBFile(t, func(t *testing.T, db KeyValueDBFile) {
		err := db.Put("key", []byte("xyz"))
		assert.NoError(t, err)
		err = db.Put("key", []byte("abc"))
		assert.NoError(t, err)
		value := db.Get("key")
		assert.Equal(t, "abc", string(value))
	})
}
