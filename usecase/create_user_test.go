package usecase

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/zackwn/article-api/repository/mongodb"
	"github.com/zackwn/article-api/security"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var createUserUseCase *CreateUserUseCase

func setup() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongo, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	db := mongo.Database("articledb_test")
	repo := mongodb.NewUserRepository(db)
	passHasher := security.NewPasswordHasher()
	createUserUseCase = NewCreateUserUseCase(repo, passHasher)
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

func TestExec(t *testing.T) {
	email := "testmail@email.com"
	err := createUserUseCase.Exec("testname", email, "d81fdw8fd81df81")
	if err != nil {
		log.Fatal(err)
	}
	_, found := createUserUseCase.userRepository.FindByEmail(email)
	if found != true {
		t.Errorf("Expect %v Got %v", "true", found)
	}
}

func TestExecErr(t *testing.T) {
	email := "testexecerr@email.com"
	createUserUseCase.Exec("testname", email, "d#8wvadw6dv7")
	err := createUserUseCase.Exec("testname", email, "dawbd8awy")
	if errors.Is(err, ErrUserAlreadyExists{Email: email}) == false {
		t.Errorf("Expect err to be ErrUserAlreadyExists Got %v", err)
	}
}
