package entity

import (
	"time"
)

func NewArticle(title, content, authorID string) (*Article, error) {
	// validation
	if title == "" {
		return nil, ArticleError{reason: "Invalid title"}
	}
	if content == "" {
		return nil, ArticleError{reason: "Invalid content"}
	}
	if authorID == "" {
		return nil, ArticleError{reason: "Invalid author_id"}
	}

	article := new(Article)
	article.ID = GenerateID()
	article.Title = title
	article.Content = content
	article.AuthorID = authorID
	article.CreatedAt = Date{time.Now()}
	return article, nil
}

type ArticleError struct {
	reason string
}

func (err ArticleError) Error() string {
	return err.reason
}

type Article struct {
	ID        string `json:"id" bson:"id"`
	Title     string `json:"title" bson:"title"`
	Content   string `json:"content" bson:"content"`
	AuthorID  string `json:"author_id" bson:"author_id"`
	CreatedAt Date   `json:"created_at" bson:"created_at"`
}
