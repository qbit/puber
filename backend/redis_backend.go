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
	var resp = true
	reply := r.client.Cmd("rpush", user, key)
	if reply.Err != nil {
		return false, reply.Err
	}
	return resp, nil
}

// RM removes a key from the lset
func (r *RedisStore) RM(user string, key string) (bool, error) {
	return true, nil
}
