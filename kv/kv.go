package kv

import (
	"github.com/rikvdh/kvui/kv/rediskv"
	"github.com/rikvdh/kvui/kv/types"
)

// KV represents a interface with functions to get and set persistent data
type KV interface {
	Databases() (int, error)
	Database(int) error
	Connected() (bool, error)
	Type(string) (types.KVType, error)

	Keys(string) ([]string, error)
	Get(string) (string, error)
	Set(string, interface{}) error
	Del(string) error

	HKeys(string) ([]string, error)
	HGet(string, string) (string, error)
	HSet(string, string, interface{}) error
	HDel(string, string) error

	LGet(string) ([]string, error)
}

const (
	// TypeRedis is a Redis KV-store
	TypeRedis string = "redis"
	// TypeRAM is a RAM-only KV-store
	TypeRAM string = "ram"
)

// New initializes a new KV-store
func New(t string, params string) (KV, error) {
	switch t {
	case TypeRedis:
		return rediskv.New(params)
		//	case TypeRAM:
		//		return ramkv.New()
	default:
		panic("invalid KV-storage")
	}
}
