package mongodb

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/zackwn/article-api/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var repo *MongoDBUserRepository

func setup() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongo, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	db := mongo.Database("articledb_test")
	repo = NewUserRepository(db)
	return db
}

func shutdown(db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := db.Collection("user").DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	db := setup()
	code := m.Run()
	shutdown(db)
	os.Exit(code)
}

func TestFindByEmailExisting(t *testing.T) {
	email := "email@mail.com"
	user, _ := entity.NewUser("name", email, "dwva8d7wafda76")
	err := repo.Store(user)
	if err != nil {
		log.Fatal(err)
	}
	result, found := repo.FindByEmail(email)
	if found != true {
		t.Errorf("Expect %v Got %v", "true", found)
	}
	if result.ID != user.ID {
		t.Errorf("Expect %v Got %v", user.ID, result.ID)
	}
}

func TestFindByEmailNotExisting(t *testing.T) {
	_, found := repo.FindByEmail("notexisting@mail.com")
	if found != false {
		t.Errorf("Expect %v Got %v", "false", found)
	}
}

func TestFindByIDExisting(t *testing.T) {
	user, _ := entity.NewUser("test", "existing@email.com", "67#afdawda8@11e")
	err := repo.Store(user)
	if err != nil {
		log.Fatal(err)
	}
	result, found := repo.FindByID(user.ID)
	if found != true {
		t.Errorf("Expect %v Got %v", "true", found)
	}
	if result.ID != user.ID {
		t.Errorf("Expect %v Got %v", user.ID, result.ID)
	}
}

func TestFindByIDNotExisting(t *testing.T) {
	_, found := repo.FindByID("notexistingID")
	if found != false {
		t.Errorf("Expect %v Got %v", "false", found)
	}
}
