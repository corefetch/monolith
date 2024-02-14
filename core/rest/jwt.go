package rest

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"time"

	"github.com/square/go-jose"
)

type Scope string

const (
	ScopeAuth  Scope = "auth"
	ScopeAdmin Scope = "admin"
)

var privateKey, _ = rsa.GenerateKey(rand.Reader, 2048)

type AuthContext struct {
	User   string    `json:"user"`
	Scope  Scope     `json:"scope"`
	Expire time.Time `json:"expire"`
}

func CreateKey(context AuthContext) (key string, err error) {

	if context.User == "" {
		return "", errors.New("expected user not empty")
	}

	if !context.Expire.After(time.Now()) {
		return "", errors.New("expected valid expiration")
	}

	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.PS512, Key: privateKey}, nil)
	if err != nil {
		return "", err
	}

	data, err := json.Marshal(context)

	if err != nil {
		return "", err
	}

	object, err := signer.Sign(data)
	if err != nil {
		return "", err
	}

	return object.CompactSerialize()
}

func ContextFromKey(key string) (context *AuthContext, err error) {

	object, err := jose.ParseSigned(key)
	if err != nil {
		return nil, err
	}

	data, err := object.Verify(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &context)

	return
}
