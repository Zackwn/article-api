package mongodb

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func shutdownTestDB(db *mongo.Database, collection string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := db.Collection(collection).DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongo, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	db := mongo.Database("articledb_test")
	setupUserTests(db)
	setupArticleTests(db)
	code := m.Run()
	shutdownTestDB(db, "user")
	shutdownTestDB(db, "article")
	os.Exit(code)
}
