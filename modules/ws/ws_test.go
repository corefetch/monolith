package ws

import (
	"sync"
	"testing"
	"time"

	"corefetch/core"
	"corefetch/core/rest"

	"github.com/gorilla/websocket"
)

func TestConnect(t *testing.T) {

	core.TestService(Service())

	key, err := rest.CreateKey(rest.AuthContext{
		User:   "1",
		Scope:  "auth",
		Expire: time.Now().Add(time.Hour),
	})

	if err != nil {
		t.Error(err)
		return
	}

	conn, _, err := websocket.DefaultDialer.Dial(
		"ws://localhost:8888/?access_token="+key,
		nil,
	)

	if err != nil {
		t.Error(err)
		return
	}

	var wg sync.WaitGroup

	wg.Add(1)

	OnMessage(func(msg Message) {

		if msg.Kind != "test" {
			t.Error("expected kind test")
		}

		if msg.Source.Context().User() != "1" {
			t.Error("expected user")
		}

		wg.Done()
	})

	if err := conn.WriteJSON(&Message{Kind: "test"}); err != nil {
		t.Error(err)
		return
	}

	wg.Wait()
}
