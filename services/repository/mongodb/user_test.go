package mongodb

import (
	"log"
	"testing"

	"github.com/zackwn/article-api/entity"
	"go.mongodb.org/mongo-driver/mongo"
)

var userRepo *MongoDBUserRepository

func setupUserTests(db *mongo.Database) {
	userRepo = NewUserRepository(db)
}

func TestUserFindByEmailExisting(t *testing.T) {
	email := "email@mail.com"
	user, _ := entity.NewUser("name", email, "picture", "dwva8d7wafda76")
	err := userRepo.Store(user)
	if err != nil {
		log.Fatal(err)
	}
	_, found := userRepo.FindByEmail(email)
	if found != true {
		t.Errorf("Expect %v Got %v", "true", found)
	}
}

func TestUserFindByEmailNotExisting(t *testing.T) {
	_, found := userRepo.FindByEmail("notexisting@mail.com")
	if found != false {
		t.Errorf("Expect %v Got %v", "false", found)
	}
}

func TestUserFindByIDExisting(t *testing.T) {
	user, _ := entity.NewUser("test", "existing@email.com", "picture", "67#afdawda8@11e")
	err := userRepo.Store(user)
	if err != nil {
		log.Fatal(err)
	}
	result, found := userRepo.FindByID(user.ID)
	if found != true {
		t.Errorf("Expect %v Got %v", "true", found)
	}
	if result.ID != user.ID {
		t.Errorf("Expect %v Got %v", user.ID, result.ID)
	}
}

func TestUserFindByIDNotExisting(t *testing.T) {
	_, found := userRepo.FindByID("notexistingID")
	if found != false {
		t.Errorf("Expect %v Got %v", "false", found)
	}
}
