package store

import (
	"context"
	"time"

	"corefetch/core/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID           primitive.ObjectID   `json:"id" bson:"_id"`
	Lesson       primitive.ObjectID   `json:"lesson"`
	Participants []primitive.ObjectID `json:"participants"`
	Duration     string               `json:"duration"`
	CreatedAt    time.Time            `json:"created_at"`
}

func (s *Session) Save() (err error) {

	_, err = db.C("vcr_sessions").InsertOne(
		context.TODO(),
		s,
	)

	return
}

func (s *Session) Drop() (err error) {

	_, err = db.C("vcr_sessions").DeleteOne(
		context.TODO(),
		bson.M{"_id": s.ID},
	)

	return
}

func GetSession(id string) (session *Session, err error) {

	_id, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	res := db.C("vcr_sessions").FindOne(
		context.TODO(),
		bson.M{"_id": _id},
	)

	err = res.Decode(&session)

	return
}
