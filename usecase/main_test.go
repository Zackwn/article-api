package usecase

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	emailservice "github.com/zackwn/article-api/services/email"
	"github.com/zackwn/article-api/services/repository/mongodb"
	"github.com/zackwn/article-api/services/security"
	temptoken "github.com/zackwn/article-api/services/temp_token"
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
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password
		DB:       0,  // default DB
	})
	ctx, cancel3 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel3()
	redisSts := redisClient.Ping(ctx)
	if redisSts.Err() != nil {
		log.Fatal("redis: ", redisSts.Err())
	}
	userRepo := mongodb.NewUserRepository(db)
	articleRepo := mongodb.NewArticleRepository(db)
	passHasher := security.NewPasswordHasher()
	emailService := emailservice.NewEmailService()
	tempToken := temptoken.NewTempToken(redisClient)
	setupCreateUserUseCase(userRepo, passHasher, emailService, tempToken)
	setupListArticlesUseCase(articleRepo, userRepo)
	code := m.Run()
	shutdownTestDB(db, "user")
	shutdownTestDB(db, "article")
	os.Exit(code)
}
