package rest

import (
	"net/http/httptest"
	"testing"
)

func TestIntent(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/?x=y", nil)
	r.Header.Set("Authorization", "AuthKey")

	i := NewContext(w, r)

	if i.Header("Authorization") != "AuthKey" {
		t.Error("expected header")
	}

	if i.Query("x") != "y" {
		t.Error("expected query")
	}
}
