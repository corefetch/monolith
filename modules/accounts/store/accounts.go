package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"learnt.io/core"
	"learnt.io/core/db"
)

type MetaKey string

type Status byte

const (
	StatusNew Status = iota + 1
	StatusChecking
	StatusApproved
	StatusRejected
)

type UserLocation struct {
	Position   *core.GeoLocation `json:"position,omitempty"`
	Country    string            `json:"country,omitempty"`
	State      string            `json:"state,omitempty"`
	City       string            `json:"city,omitempty"`
	Address    string            `json:"address,omitempty"`
	PostalCode string            `json:"postal_code,omitempty"`
}

type EmailEntry struct {
	Email    string    `json:"email"`
	Verified bool      `json:"verified,omitempty"`
	Code     *string   `json:"code" bson:"code,omitempty"`
	Created  time.Time `json:"created,omitempty"`
}

type Profile struct {
	Names                        []string  `json:"names"`
	About                        string    `json:"about,omitempty"`
	Avatar                       string    `json:"avatar,omitempty"`
	Telephone                    string    `json:"telephone,omitempty"`
	Resume                       string    `json:"resume,omitempty"`
	Birthday                     time.Time `json:"birthday,omitempty"`
	EmployerIdentificationNumber string    `json:"employer_identification_number,omitempty"`
	SocialSecurityNumber         string    `json:"social_security_number,omitempty"`
}

type Account struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Login     string             `json:"login"`
	Password  string             `json:"-"`
	Role      core.UserRole      `json:"role"`
	Profile   Profile            `json:"profile"`
	Emails    []EmailEntry       `json:"emails,omitempty"`
	Location  *UserLocation      `json:"location,omitempty"`
	Timezone  string             `json:"timezone,omitempty"`
	Status    Status             `json:"status,omitempty" bson:"status,omitempty"`
	Meta      map[MetaKey]any    `json:"meta,omitempty" bson:"meta,omitempty"`
	CreatedAt time.Time          `json:"created_at"`
}

func (a *Account) Email() string {
	return a.Login
}

func (a *Account) Name() string {
	return strings.Join(a.Profile.Names, " ")
}

func (a *Account) HasMeta(key MetaKey) bool {
	_, has := a.Meta[key]
	return has
}

func (a *Account) GetMeta(key MetaKey) any {
	v, has := a.Meta[key]
	if !has {
		panic("expected key")
	}
	return v
}

// Remove an account
func (a *Account) Drop() (err error) {

	_, err = db.C("accounts").DeleteOne(
		context.TODO(),
		bson.M{"_id": a.ID},
	)

	return
}

func (user *Account) Save() (err error) {

	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}

	_, err = db.C("accounts").InsertOne(
		context.TODO(),
		user,
	)
	return
}

func GetAccountByLogin(login string) (user *Account, err error) {

	res := db.C("accounts").FindOne(
		context.TODO(),
		bson.M{"login": login},
	)

	err = res.Decode(&user)

	if err != nil {
		return nil, errors.New("no user found")
	}

	return
}

func GetAccount(id string) (user *Account, err error) {

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	res := db.C("accounts").FindOne(
		context.TODO(),
		bson.M{"_id": oid},
	)

	err = res.Decode(&user)

	return
}
