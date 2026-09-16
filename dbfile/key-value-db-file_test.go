package dbfile

import (
	"fmt"
	"math/rand"
	"os"
	"slices"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func withKeyValueDBFile(t *testing.T, f func(t *testing.T, db KeyValueDBFile)) {
	fn := fmt.Sprintf("test-%d.tmp", rand.Intn(1000))
	db, err := NewKeyValueDBFile(fn)
	assert.NoError(t, err)
	f(t, db)
	_ = os.Remove(fn)
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

// Meaningful under -race.
func TestConcurrentAccess(t *testing.T) {
	withKeyValueDBFile(t, func(t *testing.T, db KeyValueDBFile) {
		var wg sync.WaitGroup
		for i := 0; i < 4; i++ {
			wg.Add(2)
			go func(i int) {
				defer wg.Done()
				for j := 0; j < 10; j++ {
					assert.NoError(t, db.Put(fmt.Sprintf("k%d", i), []byte("v")))
				}
			}(i)
			go func(i int) {
				defer wg.Done()
				for j := 0; j < 10; j++ {
					_ = db.Get(fmt.Sprintf("k%d", i))
					_ = db.Keys()
				}
			}(i)
		}
		wg.Wait()
	})
}

func TestFileIsNotExecutable(t *testing.T) {
	withKeyValueDBFile(t, func(t *testing.T, db KeyValueDBFile) {
		assert.NoError(t, db.Put("key", []byte("xyz")))
		info, err := os.Stat(db.(*keyFile).file)
		assert.NoError(t, err)
		assert.Equal(t, os.FileMode(0), info.Mode().Perm()&0o133, "file must not be executable or group/world writable")
	})
}

func TestPersistedAcrossReopen(t *testing.T) {
	withKeyValueDBFile(t, func(t *testing.T, db KeyValueDBFile) {
		fn := db.(*keyFile).file
		assert.NoError(t, db.Put("key", []byte("xyz")))
		_, err := os.Stat(fn + ".tmp")
		assert.True(t, os.IsNotExist(err), "temporary file must not be left behind")
		reopened, err := NewKeyValueDBFile(fn)
		assert.NoError(t, err)
		assert.Equal(t, "xyz", string(reopened.Get("key")))
	})
}
