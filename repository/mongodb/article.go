package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zackwn/article-api/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewArticleRepository(mongodb *mongo.Database) *ArticleRepository {
	return &ArticleRepository{db: mongodb}
}

type ArticleRepository struct {
	db *mongo.Database
}

func (repo ArticleRepository) FindByID(ID string) (*entity.Article, bool) {
	collection := repo.db.Collection("article")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	filter := bson.D{{Key: "id", Value: ID}}
	defer cancel()
	var result entity.Article
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Fatal(err)
		}
		return nil, false
	}
	return &result, true
}

func (repo ArticleRepository) FindAllByAuthor(authorID string) []*entity.Article {
	collection := repo.db.Collection("article")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	filter := bson.D{{Key: "author_id", Value: authorID}}
	defer cancel()
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	var results []*entity.Article
	for cursor.Next(ctx) {
		var result entity.Article
		cursor.Decode(&result)
		fmt.Println("here")
		results = append(results, &result)
	}
	return results
}

func (repo ArticleRepository) Store(article *entity.Article) error {
	collection := repo.db.Collection("article")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, article)
	if err != nil {
		return err
	}
	return nil
}
