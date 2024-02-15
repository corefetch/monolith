package ws

import (
	"net/http"
	"sync"
	"testing"
	"time"

	"corefetch/core/rest"

	"github.com/gorilla/websocket"
	"github.com/olebedev/emitter"
)

func TestConnect(t *testing.T) {

	go http.ListenAndServe(":5678", Service())

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
		"ws://localhost:5678/?access_token="+key,
		nil,
	)

	if err != nil {
		t.Error(err)
		return
	}

	var wg sync.WaitGroup

	wg.Add(1)

	events.On("message", func(e *emitter.Event) {
		msg := e.Args[0].(*Message)
		if msg.Kind != "test" {
			t.Error("expected kind test")
		}
		wg.Done()
	})

	if err := conn.WriteJSON(&Message{Kind: "test"}); err != nil {
		t.Error(err)
		return
	}

	wg.Wait()
}
