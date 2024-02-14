package rest

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCases(t *testing.T) {

	tests := []map[string]any{
		{
			"context": AuthContext{
				User:   "1",
				Scope:  ScopeAuth,
				Expire: time.Now().Add(time.Hour),
			},
			"scope": ScopeAuth,
			"pass":  true,
		},
		{
			"context": AuthContext{
				User:   "1",
				Scope:  ScopeAuth,
				Expire: time.Now().Add(time.Hour),
			},
			"scope": ScopeAdmin,
			"pass":  false,
		},
		{
			"context": AuthContext{
				User:   "1",
				Scope:  ScopeAuth,
				Expire: time.Now().Add(time.Hour),
			},
			"scope": ScopeAdmin,
			"pass":  false,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Test Auth %d", i), func(t *testing.T) {

			key, err := CreateKey(test["context"].(AuthContext))

			if err != nil {
				t.Error(err)
				return
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)
			r.Header.Add("Authorization", "Bearer "+key)

			var passed bool

			GuardScope(test["scope"].(Scope), func(i *Context) {
				passed = true
			})(NewContext(w, r))

			if passed != test["pass"].(bool) {
				t.Error("pass expected")
			}
		})
	}
}

func TestNoKey(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	var passed bool

	GuardAuth(func(i *Context) {
		passed = true
	})(NewContext(w, r))

	if passed {
		t.Error("pass expected not")
	}
}
