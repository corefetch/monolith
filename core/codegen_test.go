package core

import "testing"

func TestShortCodeGen(t *testing.T) {
	key := NewShortKey()
	if len(key) != 6 {
		t.Error("expected key length of 6")
	}
}
