package db

import (
	"context"
	"net/url"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func Session() *mongo.Database {

	if db != nil {
		return db
	}

	if os.Getenv("DB") == "" {
		panic("DB env uri expected")
	}

	config := options.Client().ApplyURI(
		os.Getenv("DB"),
	)

	client, err := mongo.Connect(context.TODO(), config)

	if err != nil {
		panic(err)
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		panic(err)
	}

	u, err := url.Parse(os.Getenv("DB"))

	if err != nil {
		panic(err)
	}

	database, _ := strings.CutPrefix(u.Path, "/")

	db = client.Database(database)

	return db
}

// Get collection by name
func C(name string) *mongo.Collection {
	return Session().Collection(name)
}
