package ws

import (
	"net/http"
	"testing"
	"time"

	"corefetch/core/rest"

	"github.com/gorilla/websocket"
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

	if err := conn.WriteJSON(&Message{Kind: "ping"}); err != nil {
		t.Error(err)
		return
	}

	var message Message

	if err := conn.ReadJSON(&message); err != nil {
		t.Error(err)
		return
	}

	if message.Kind != "pong" {
		t.Error("expected pong")
	}
}
