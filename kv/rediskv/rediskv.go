package rediskv

import (
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/rikvdh/redisui/kv/types"
)

type redisCon interface {
	Err() error
	Do(cmd string, args ...interface{}) (interface{}, error)
}

// Rediskv stores the values that is set or retrieved in Redis
type Rediskv struct {
	redis redisCon
}

// Get returns the value from the requested key.
func (r Rediskv) Get(key string) (string, error) {
	v, e := r.redis.Do("GET", key)
	log.Println(v)
	return redis.String(v, e)
}

// Set stores the value with the given key
func (r Rediskv) Set(key string, value interface{}) error {
	_, err := r.redis.Do("SET", key, value)
	return err
}

// Keys returns a list of keys matched by pattern
func (r Rediskv) Keys(pattern string) ([]string, error) {
	// Note: maybe use scan instead of keys. Keys has a great performance impact
	// See: https://redis.io/commands/keys, https://redis.io/commands/scan
	return redis.Strings(r.redis.Do("KEYS", pattern))
}

// Del removes the value with the given key
func (r Rediskv) Del(key string) error {
	_, err := r.redis.Do("DEL", key)
	return err
}

// HKeys returns all keys in a hash
func (r Rediskv) HKeys(key string) ([]string, error) {
	return redis.Strings(r.redis.Do("HKEYS", key))
}

// HGet retrieve the value from the given field in the given key
func (r Rediskv) HGet(key, field string) (string, error) {
	return redis.String(r.redis.Do("HGET", key, field))
}

// HSet sets the value in the given field in the given key
func (r Rediskv) HSet(key, field string, value interface{}) error {
	_, err := r.redis.Do("HSET", key, field, value)
	return err
}

// HDel removes the value in the given field in the given key
func (r Rediskv) HDel(key, field string) error {
	_, err := r.redis.Do("HDEL", key, field)
	return err
}

// Databases requested from config
func (r Rediskv) Databases() (int, error) {
	ret, err := redis.Values(r.redis.Do("CONFIG", "GET", "databases"))
	if len(ret) > 1 {
		return redis.Int(ret[1], err)
	}
	return 0, err
}

func (r Rediskv) Connected() (bool, error) {
	err := r.redis.Err()
	ret := true
	if err != nil {
		ret = false
	}
	return ret, err
}

func (r Rediskv) Type(key string) (types.KVType, error) {
	t, err := redis.String(r.redis.Do("TYPE", key))
	if err == nil {
		return r.redisTypeToKVType(t)
	}
	return types.KVTypeInvalid, err
}

func (r Rediskv) redisTypeToKVType(t string) (types.KVType, error) {
	switch t {
	case "hash":
		return types.KVTypeMap, nil
	case "string":
		return types.KVTypeString, nil
	case "list":
		return types.KVTypeList, nil
	}
	return types.KVTypeInvalid, fmt.Errorf("invalid type: %s", t)
}

// New creates a Redis key value instance
func New(host string) (*Rediskv, error) {
	rediskv := Rediskv{}
	redisCon, err := redis.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	rediskv.redis = redisCon

	return &rediskv, nil
}
