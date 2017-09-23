package kv

import (
	"net"
	"testing"

	"github.com/rikvdh/kvui/kv/ramkv"
	"github.com/rikvdh/kvui/kv/rediskv"
)

func TestKvRedis(t *testing.T) {
	l, err := net.Listen("tcp", ":56789")
	if err != nil {
		panic(err)
	}
	defer l.Close()

	r, err := New(TypeRedis, "127.0.0.1:56789")
	if err != nil {
		t.Error("Error must not be nil")
	}
	if _, ok := r.(*rediskv.Rediskv); !ok {
		t.Error("Reply must be a rediskv")
	}
}

func TestKvRAM(t *testing.T) {
	r, err := New(TypeRAM, "")
	if err != nil {
		t.Error("Error must not be nil")
	}
	if _, ok := r.(*ramkv.Ramkv); !ok {
		t.Error("Reply must be a rediskv")
	}
}

func TestInvalidKV(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("The code dit not panic")
		}
	}()
	New("boem", "")
}
