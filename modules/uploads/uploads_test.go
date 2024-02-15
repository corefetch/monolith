package uploads

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"testing"
	"time"

	"corefetch/core"
	"corefetch/core/rest"
	"corefetch/modules/accounts/store"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUploadAndRender(t *testing.T) {

	core.TestService(Service())

	var userID = primitive.NewObjectID()

	account := store.Account{
		ID: userID,
	}

	account.Save()
	defer account.Drop()

	var buf = &bytes.Buffer{}

	mw := multipart.NewWriter(buf)

	w, err := mw.CreateFormFile("upload", "uploads.txt")

	if err != nil {
		t.Error(err)
		return
	}

	if _, err := w.Write([]byte("Hello World")); err != nil {
		t.Error(err)
		return
	}

	mw.Close()

	req, err := http.NewRequest("POST", "http://localhost:8888/", buf)

	if err != nil {
		t.Error(err)
		return
	}

	key, err := rest.CreateKey(rest.AuthContext{
		User:   userID.Hex(),
		Scope:  "auth",
		Expire: time.Now().Add(time.Minute),
	})

	if err != nil {
		t.Error(err)
		return
	}

	req.Header.Add("Content-Type", mw.FormDataContentType())
	req.Header.Add("Authorization", "Bearer "+key)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Error(err)
		return
	}

	if res.StatusCode != 201 {
		t.Error("expected 201:", res.StatusCode)
		return
	}

	data, err := io.ReadAll(res.Body)

	if err != nil {
		t.Error(err)
		return
	}

	var upload Upload

	if err := json.Unmarshal(data, &upload); err != nil {
		t.Error(err)
		return
	}

	req, err = http.NewRequest("GET", "http://localhost:8888/"+upload.ID.Hex(), nil)
	req.Header.Add("Authorization", "Bearer "+key)

	if err != nil {
		t.Error(err)
		return
	}

	res, err = http.DefaultClient.Do(req)

	if err != nil {
		t.Error(err)
		return
	}

	data, err = io.ReadAll(res.Body)

	if err != nil {
		t.Error(err)
		return
	}

	if string(data) != "Hello World" {
		t.Error("expected render text but:", string(data))
		return
	}

	req, err = http.NewRequest("DELETE", "http://localhost:8888/"+upload.ID.Hex(), nil)
	req.Header.Add("Authorization", "Bearer "+key)

	if err != nil {
		t.Error(err)
		return
	}

	if res.StatusCode != 200 {
		t.Error("expected delete")
		return
	}
}
