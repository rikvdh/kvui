package rediskv

import (
	"fmt"
	"testing"
)

type redisMock struct {
	Result interface{}
	err    error
}

func (r redisMock) Do(cmd string, args ...interface{}) (interface{}, error) {
	return r.Result, r.err
}

func (r redisMock) Err() error {
	return r.err
}

func TestGetSet(t *testing.T) {
	mock := redisMock{}
	kvStorage := Rediskv{}
	kvStorage.redis = &mock

	mock.Result = nil
	mock.err = nil

	err := kvStorage.Set("value", `"{"name": "simulator", "value": "20"}"`)
	if err != nil {
		t.Error("Couln't set value name:simulator, value:20")
	}

	mock.Result = `"{"name": "simulator", "value": "20"}"`
	mock.err = nil

	value, err := kvStorage.Get("value")
	if err != nil {
		t.Error("Couln't get value for initialized value")
		t.Error(err)
	}

	if value != `"{"name": "simulator", "value": "20"}"` {
		t.Error("Unexpected get value for initialized value")
	}
}

func TestDel(t *testing.T) {
	mock := redisMock{}
	kvStorage := Rediskv{}
	kvStorage.redis = &mock

	mock.Result = nil
	mock.err = nil

	err := kvStorage.Set("value", `{"name": "simulator", "value": "20"}`)
	if err != nil {
		t.Error("Couln't set value name:simulator, value:20")
	}

	mock.Result = `{"name": "simulator", "value": "20"}`
	mock.err = nil

	value, err := kvStorage.Get("value")
	if err != nil {
		t.Error("Couln't get value for initialized value")
	}

	if value != `{"name": "simulator", "value": "20"}` {
		t.Errorf("Unexpected get value for initialized value got: %v", value)
	}

	mock.Result = nil
	mock.err = nil

	err = kvStorage.Del("value")
	if err != nil {
		t.Error("Couln't del value")
	}

	mock.Result = nil
	mock.err = nil

	value, err = kvStorage.Get("value")
	if err == nil {
		t.Error("Couln't get value for initialized value")
	}

	if value != "" {
		t.Error("Unexpected get value for removed value")
	}
}

func TestHGetHSet(t *testing.T) {
	mock := redisMock{}
	kvStorage := Rediskv{}
	kvStorage.redis = &mock

	mock.Result = nil
	mock.err = nil

	err := kvStorage.HSet("value", "name", "simulator")
	if err != nil {
		t.Error("Couln't hset value with field name and value simulator")
	}

	mock.Result = "simulator"
	mock.err = nil

	value, err := kvStorage.HGet("value", "name")
	if err != nil {
		t.Error("Couln't hget value for field name")
	}

	if value != "simulator" {
		t.Error("Unexpected hget value for field name")
	}
}

func TestHDel(t *testing.T) {
	mock := redisMock{}
	kvStorage := Rediskv{}
	kvStorage.redis = &mock

	mock.Result = nil
	mock.err = nil

	err := kvStorage.HSet("value", "name", "simulator")
	if err != nil {
		t.Error("Couln't hset value with field name and value simulator")
	}

	mock.Result = "simulator"
	mock.err = nil

	value, err := kvStorage.HGet("value", "name")
	if err != nil {
		t.Error("Couln't hget value for field name")
	}

	if value != "simulator" {
		t.Error("Unexpected hget value for field name")
	}

	mock.Result = nil
	mock.err = nil

	err = kvStorage.HDel("value", "name")
	if err != nil {
		t.Error("Couln't del value")
	}
}

func TestKeys(t *testing.T) {
	mock := redisMock{}
	kvStorage := Rediskv{}
	kvStorage.redis = &mock

	mock.Result = nil
	mock.err = fmt.Errorf("test-err")

	k, err := kvStorage.Keys("*")
	if err.Error() != "test-err" {
		t.Errorf("keys error didnt match: %s (%v)", err.Error(), err)
	}
	if k != nil {
		t.Error("keys must be nil")
	}
}
