package mongodb

import (
	"context"
	"log"
	"time"

	"github.com/zackwn/article-api/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewUserRepository(mongodb *mongo.Database) *MongoDBUserRepository {
	return &MongoDBUserRepository{db: mongodb}
}

type MongoDBUserRepository struct {
	db *mongo.Database
}

// FindByEmail returns (user, found)
func (repo MongoDBUserRepository) FindByEmail(email string) (*entity.User, bool) {
	collection := repo.db.Collection("user")
	var result entity.User
	filter := bson.D{{Key: "email", Value: email}}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Println(err)
		}
		return nil, false
	}
	return &result, true
}

// FindByID returns (user, found)
func (repo MongoDBUserRepository) FindByID(ID string) (*entity.User, bool) {
	collection := repo.db.Collection("user")
	var result entity.User
	filter := bson.D{{Key: "id", Value: ID}}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Println(err)
		}
		return nil, false
	}
	return &result, true
}

func (repo MongoDBUserRepository) Store(user *entity.User) error {
	collection := repo.db.Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (repo MongoDBUserRepository) Update(user *entity.User) error {
	collection := repo.db.Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{{Key: "id", Value: user.ID}}
	update := bson.M{
		"$set": user,
	}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}
