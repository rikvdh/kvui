package ramkv

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Ramkv stores the values that is set or retrieved in RAM
type Ramkv struct {
	lock    sync.RWMutex
	storage map[string]map[string]string
}

// Get returns the value from the requested key.
func (r *Ramkv) Get(key string) (string, error) {
	r.lock.RLock()
	value, found := r.storage[key]
	r.lock.RUnlock()
	if found && value != nil {
		value, err := json.Marshal(value)
		if err != nil {
			return "", err
		}
		return string(value), nil
	}
	return "", fmt.Errorf("key %s not found", key)
}

// Set stores the value with the given key
func (r *Ramkv) Set(key string, value interface{}) error {
	r.lock.Lock()
	if r.storage[key] == nil {
		r.storage[key] = make(map[string]string)
	}
	r.lock.Unlock()

	if strValue, ok := value.(string); ok {
		var jsonValue map[string]string
		err := json.Unmarshal([]byte(strValue), &jsonValue)
		if err != nil {
			return err
		}
		r.lock.Lock()
		r.storage[key] = jsonValue
		r.lock.Unlock()
	}

	return nil
}

// Del removes the value with the given key
func (r *Ramkv) Del(key string) error {
	r.lock.Lock()
	delete(r.storage, key)
	r.lock.Unlock()
	return nil
}

// Keys returns a list of all keys, pattern is ignored
func (r *Ramkv) Keys(pattern string) ([]string, error) {
	var keys []string

	r.lock.RLock()
	for k := range r.storage {
		keys = append(keys, k)
	}
	r.lock.RUnlock()

	return keys, nil
}

// HKeys returns all keys in a hash
func (r *Ramkv) HKeys(key string) ([]string, error) {
	var keys []string

	r.lock.RLock()
	s, ok := r.storage[key]
	r.lock.RUnlock()
	if !ok {
		return nil, fmt.Errorf("key %s not found", key)
	}

	for k := range s {
		keys = append(keys, k)
	}

	return keys, nil
}

// HGet retrieve the value from the given field in the given key
func (r *Ramkv) HGet(key, field string) (string, error) {
	r.lock.RLock()
	if r.storage[key] == nil {
		r.lock.RUnlock()
		return "", fmt.Errorf("key %s not found", key)
	}

	value, found := r.storage[key][field]
	r.lock.RUnlock()
	if found {
		return value, nil
	}
	return "", fmt.Errorf("field %s not found", field)
}

// HSet sets the value in the given field in the given key
func (r *Ramkv) HSet(key, field string, value interface{}) error {
	r.lock.Lock()
	if r.storage[key] == nil {
		r.storage[key] = make(map[string]string)
	}

	r.storage[key][field] = fmt.Sprintf("%v", value)
	r.lock.Unlock()
	return nil
}

// HDel removes the value in the given field in the given key
func (r *Ramkv) HDel(key, field string) error {
	r.lock.Lock()
	if r.storage[key] != nil {
		delete(r.storage[key], field)
	}
	r.lock.Unlock()
	return nil
}

// Databases is always 1 for RAM
func (*Ramkv) Databases() (int, error) {
	return 1, nil
}

// Connected is always true for RAM
func (*Ramkv) Connected() (bool, error) {
	return true, nil
}

// New creates a Ram key value instance
func New() (*Ramkv, error) {
	ramkv := Ramkv{}
	ramkv.storage = make(map[string]map[string]string)
	return &ramkv, nil
}
