package rest

import (
	"testing"
	"time"
)

func TestKeyEncoder(t *testing.T) {

	context := AuthContext{
		User:   "1",
		Scope:  "auth",
		Expire: time.Now().Add(time.Minute * 5),
	}

	data, err := CreateKey(context)

	if err != nil {
		t.Error(err)
		return
	}

	out, err := ContextFromKey(data)

	if err != nil {
		t.Error(err)
		return
	}

	if out.User != "1" {
		t.Error("expected user match")
	}

	if out.Scope != "auth" {
		t.Error("expected user match")
	}
}
