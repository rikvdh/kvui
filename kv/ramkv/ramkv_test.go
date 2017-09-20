package ramkv

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSet(t *testing.T) {
	kvStorage, err := New()
	assert.Nil(t, err)

	value, err := kvStorage.Get("value")
	if err.Error() != "key value not found" {
		t.Error("Expected error for ramkv.Get")
	}

	if value != "" {
		t.Error("Unexpected get value for uninitialized value")
	}

	err = kvStorage.Set("value", `{"name": "simulator", "value": "20"}`)
	if err != nil {
		t.Error("Couln't set value name:simulator, value:20")
	}

	value, err = kvStorage.Get("value")
	if err != nil {
		t.Error("Couln't get value for initialized value")
	}

	if value != `{"name":"simulator","value":"20"}` {
		t.Error("Unexpected get value for initialized value")
	}
}

func TestDel(t *testing.T) {
	kvStorage, err := New()
	assert.Nil(t, err)

	err = kvStorage.Set("value", `{"name": "simulator", "value": "20"}`)
	assert.Nil(t, err)

	value, err := kvStorage.Get("value")
	assert.Nil(t, err)

	if value != `{"name":"simulator","value":"20"}` {
		t.Error("Unexpected get value for value")
	}

	err = kvStorage.Del("value")
	assert.Nil(t, err)

	value, err = kvStorage.Get("value")
	if err == nil {
		t.Error("error expected, no value should be present")
	}

	if value != "" {
		t.Error("Unexpected get value for removed value")
	}
}

func TestHGetHSet(t *testing.T) {
	kvStorage, err := New()
	assert.Nil(t, err)

	value, err := kvStorage.HGet("value", "name")
	if err.Error() != "key value not found" {
		t.Error("Expected error for ramkv.Get")
	}

	if value != "" {
		t.Error("Unexpected hget value for uninitialized value")
	}

	err = kvStorage.HSet("value", "name", "simulator")
	assert.Nil(t, err)

	value, err = kvStorage.HGet("value", "name")
	assert.Nil(t, err)

	if value != "simulator" {
		t.Error("Unexpected hget value for field name")
	}

	err = kvStorage.HSet("value", "name", 100)
	if err != nil {
		t.Error("Couln't hset value with field name and value simulator")
	}

	value, err = kvStorage.HGet("value", "name")
	assert.Nil(t, err)

	if value != "100" {
		t.Error("Unexpected hget value for field name")
	}

	value, err = kvStorage.HGet("value", "value")
	if err.Error() != "field value not found" {
		t.Error("Expected error for ramkv.Get")
	}

	if value != "" {
		t.Error("Unexpected hget value for field value")
	}
}

func TestHDelNotExisting(t *testing.T) {
	kvStorage, err := New()
	assert.Nil(t, err)

	err = kvStorage.HSet("value", "name", "boem")
	if err != nil {
		t.Error("no error expected for HSet")
	}

	err = kvStorage.HDel("value", "rik")
	assert.Nil(t, err)

	err = kvStorage.HDel("friet", "boembats")
	assert.Nil(t, err)
}

func TestHDel(t *testing.T) {
	kvStorage, err := New()
	assert.Nil(t, err)

	err = kvStorage.HSet("value", "name", "simulator")
	assert.Nil(t, err)

	value, err := kvStorage.HGet("value", "name")
	assert.Nil(t, err)

	if value != "simulator" {
		t.Error("Unexpected hget value for field name")
	}

	err = kvStorage.HDel("value", "name")
	assert.Nil(t, err)

	value, err = kvStorage.HGet("value", "name")
	if err == nil {
		t.Errorf("hget should fail")
	}

	if value != "" {
		t.Error("Unexpected hget value for removed value")
	}
}

func TestKeys(t *testing.T) {
	kvStorage, err := New()
	assert.Nil(t, err)

	k, err := kvStorage.Keys("*")
	assert.Nil(t, err)

	if k != nil {
		t.Error("keys must be nil")
	}

	kvStorage.Set("whots", "bats")
	kvStorage.HSet("knal", "boem", "beng")

	k, err = kvStorage.Keys("*")
	assert.Nil(t, err)

	if k == nil {
		t.Error("keys must not be nil")
	}
	if len(k) != 2 {
		t.Errorf("expected 2 keys")
	}
}