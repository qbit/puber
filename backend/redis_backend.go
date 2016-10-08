package backend

import (
	"time"

	"github.com/fzzy/radix/redis"
)

// RedisStore structure for our redis store
type RedisStore struct {
	client *redis.Client
}

// Init connects to our local redis store (assumption is it's running defaults)
func (r *RedisStore) Init() error {
	var err error
	r.client, err = redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
	return err
}

// Close ends our session with redis
func (r *RedisStore) Close() {
	r.client.Close()
}

// Add inserts a key into the redis backend
func (r *RedisStore) Add(user string, key string) (bool, error) {
	reply := r.client.Cmd("rpush", user, key)
	if reply.Err != nil {
		return false, reply.Err
	}
	return true, nil
}

// RM removes a key from the lset
func (r *RedisStore) RM(user string, key string) (bool, error) {
	reply := r.client.Cmd("lrem", user, -1, key)

	if reply.Err != nil {
		return false, reply.Err
	}
	return true, nil
}

// RMAll removes all the keys for a given user
func (r *RedisStore) RMAll(user string) (bool, error) {
	reply := r.client.Cmd("del", user)
	if reply.Err != nil {
		return false, reply.Err
	}

	return true, nil
}

// Get returns all keys for user
func (r *RedisStore) Get(user string) ([]string, error) {
	keys, err := r.client.Cmd("lrange", user, 0, 1).List()
	if err != nil {
		return nil, err
	}

	return keys, nil
}

// GetAll returns all keys for all users
func (r *RedisStore) GetAll() ([]string, error) {
	var s []string
	reply := r.client.Cmd("keys", "*")
	keys, err := reply.List()
	if err != nil {
		return nil, err
	}

	for _, k := range keys {
		uKeys, err := r.client.Cmd("lrange", k, 0, -1).List()
		if err != nil {
			return nil, err
		}
		for i := range uKeys {
			s = append(s, uKeys[i])
		}
	}

	return s, nil
}

// GetCount gets a count of all the keys
func (r *RedisStore) GetCount() (int, error) {
	s := 0
	keys, err := r.client.Cmd("keys", "*").List()
	if err != nil {
		return -1, err
	}

	for _, k := range keys {
		uKeys, err := r.client.Cmd("lrange", k, 0, -1).List()
		if err != nil {
			return -1, err
		}
		for _ = range uKeys {
			s++
		}
	}

	return s, nil
}
