package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	c "github.com/zackwn/article-api/controllers"
	"github.com/zackwn/article-api/repository/mongodb"
	"github.com/zackwn/article-api/security"
	"github.com/zackwn/article-api/usecase"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func adaptController(method string, controller c.Controller) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if req.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		res := controller.Handle(req)
		w.WriteHeader(res.StatusCode)
		if len(res.Data) != 0 {
			w.Write(res.Data)
		}
	}
}

func main() {
	fmt.Println("Connecting to the Database...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel2 := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel2()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database("articledb")
	fmt.Println("Successfully connected to the Database.")

	// services
	userRepository := mongodb.NewUserRepository(db)
	articleRepository := mongodb.NewArticleRepository(db)
	passwordHasher := security.NewPasswordHasher()
	authProvider := security.NewAuthProvider()

	// usecases
	createUserUseCase := usecase.NewCreateUserUseCase(userRepository, passwordHasher)
	userLoginUseCase := usecase.NewUserLoginUseCase(userRepository, passwordHasher, authProvider)
	createArticleUseCase := usecase.NewCreateArticleUseCase(authProvider, articleRepository, userRepository)
	listArticles := usecase.NewListArticlesUseCase(articleRepository, userRepository)

	// controllers
	userSignupController := c.NewUserSignupController(createUserUseCase)
	userSigninController := c.NewUserSigninController(userLoginUseCase)
	createArticleController := c.NewCreateArticleController(createArticleUseCase)
	listArticlesController := c.NewListArticlesController(listArticles)

	http.HandleFunc("/user/signup", adaptController("POST", userSignupController))
	http.HandleFunc("/user/signin", adaptController("POST", userSigninController))
	http.HandleFunc("/articles/create", adaptController("POST", createArticleController))
	http.HandleFunc("/articles/list", adaptController("GET", listArticlesController))

	fmt.Println("Server started...")
	http.ListenAndServe(":8080", nil)
}
