/*
 * @Date: 2024-06-20 09:48:51
 * @LastEditTime: 2024-06-21 09:48:52
 * @Description:
 */
package memory

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/coocood/freecache"
)

type Store struct {
	store *freecache.Cache
}

// New
// @param conf
// @date 2022-07-02 08:12:14
func New(cacheSize int) *Store {
	if cacheSize == 0 {
		cacheSize = 100 * 1024 * 1024
	}
	cache := freecache.NewCache(cacheSize)
	return &Store{store: cache}
}

// NewWithDb
// @param tx
// @date 2022-07-02 08:12:12
func NewWithDb(tx *freecache.Cache) *Store {
	return &Store{store: tx}
}

// Set
// @param ctx
// @param key
// @param value
// @param ttl
// @date 2022-07-02 08:12:11
func (r *Store) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	seconds := ttl.Seconds()
	valType := reflect.TypeOf(value)
	bytes := []byte{}
	switch valType {
	case reflect.TypeOf([]byte{}):
		bytes = []byte(value.([]byte))
	case reflect.TypeOf(string("")):
		bytes = []byte(value.(string))
	case reflect.TypeOf([]uint8{}):
		bs := value.([]uint8)
		for _, b := range bs {
			bytes = append(bytes, byte(b))
		}
	}
	return r.store.Set([]byte(key), bytes, int(seconds))
}

// Get
// @param ctx
// @param key
// @date 2022-07-02 08:12:09
func (r *Store) Get(ctx context.Context, key string) ([]byte, error) {
	return r.store.Get([]byte(key))
}

// RemoveFromTag
// @param ctx
// @param tag
// @date 2022-07-02 08:12:08
func (r *Store) RemoveFromTag(ctx context.Context, tag string) error {
	keys, err := r.store.Get([]byte(tag))
	if err != nil {
		return err
	}
	keysStr := string(keys)
	if len(keysStr) > 0 {
		keyList := strings.Split(keysStr, ",")
		for _, k := range keyList {
			r.store.Del([]byte(k))
		}
	}
	return nil
}

func (r *Store) RemoveFromKey(ctx context.Context, key string) error {
	affected := r.store.Del([]byte(key))
	if affected {
		return nil
	} else {
		return errors.New("key not found")
	}
}

// SaveTagKey
// @param ctx
// @param tag
// @param key
// @date 2022-07-02 08:12:05
func (r *Store) SaveTagKey(ctx context.Context, tag, key string) error {
	keys, _ := r.store.Get([]byte(tag))
	if len(keys) == 0 {
		keys = []byte(key)
	} else {
		keysStr := string(keys)
		keyList := strings.Split(keysStr, ",")
		found := false
		for _, k := range keyList {
			if k == key {
				found = true
				break
			}
		}
		if found {
			return nil
		} else {
			keysStr = keysStr + "," + key
		}
		keys = []byte(keysStr)
	}
	return r.store.Set([]byte(tag), keys, 0)
}

// RemoveTagKey
// @param ctx
// @param tag
// @param key
// @date 2022-07-02 08:12:05
func (r *Store) RemoveTagKey(ctx context.Context, tag, key string) error {
	keys, err := r.store.Get([]byte(tag))
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	keysStr := string(keys)
	keyList := strings.Split(keysStr, ",")
	var result []string
	for _, k := range keyList {
		if k != key {
			result = append(result, k)
		}
	}
	return r.store.Set([]byte(tag), []byte(strings.Join(result, ",")), 0)
}

// MemberTagKey
// @param ctx
// @param tag
// @param key
// @date 2022-07-02 08:12:05
func (r *Store) MemberTagKey(ctx context.Context, tag, key string) (bool, error) {
	keys, err := r.store.Get([]byte(tag))
	if err != nil {
		return false, err
	}
	if len(keys) == 0 {
		return false, nil
	}
	keysStr := string(keys)
	keyList := strings.Split(keysStr, ",")
	found := false
	for _, k := range keyList {
		if k == key {
			found = true
			break
		}
	}
	return found, nil
}

func (r *Store) Clear(ctx context.Context, keys ...string) error {
	r.store.Clear()
	return nil
}
