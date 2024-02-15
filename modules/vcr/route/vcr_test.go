package route

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"corefetch/core"
	"corefetch/core/rest"
	"corefetch/modules/vcr/store"
	"corefetch/modules/ws"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestJoin(t *testing.T) {

	core.TestSetup()

	sessID := primitive.NewObjectID()

	session := store.Session{
		ID:           sessID,
		Lesson:       primitive.NewObjectID(),
		Participants: make([]primitive.ObjectID, 0),
		Duration:     "1h",
		CreatedAt:    time.Now(),
	}

	session.Save()
	defer session.Drop()

	connect(t)

	data := bytes.NewBufferString(fmt.Sprintf(`{"session":"%s"}`, sessID.Hex()))

	key, err := rest.CreateKey(rest.AuthContext{
		User:   "1",
		Scope:  "auth",
		Expire: time.Now().Add(time.Hour),
	})

	if err != nil {
		t.Error("not able to create auth key")
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", data)
	r.Header.Add("Authorization", "Bearer "+key)
	ctx := rest.NewContext(w, r)

	rest.GuardAuth(Join)(ctx)

	if w.Result().StatusCode != 200 {
		t.Error("expected 200:", w.Result().StatusCode)
	}

	_, err = ws.UserConnection("1")

	if err != nil {
		t.Error(err)
		return
	}
}

func connect(t *testing.T) {

	go http.ListenAndServe(":5678", ws.Service())

	key, err := rest.CreateKey(rest.AuthContext{
		User:   "1",
		Scope:  "auth",
		Expire: time.Now().Add(time.Hour),
	})

	if err != nil {
		t.Error(err)
		return
	}

	_, _, err = websocket.DefaultDialer.Dial(
		"ws://localhost:5678/?access_token="+key,
		nil,
	)

	if err != nil {
		t.Error(err)
	}
}
