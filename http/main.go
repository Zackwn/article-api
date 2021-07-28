package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	c "github.com/zackwn/article-api/controllers"
	emailservice "github.com/zackwn/article-api/services/email"
	filestorage "github.com/zackwn/article-api/services/file_storage"
	fph "github.com/zackwn/article-api/services/forgot_password_handler"
	"github.com/zackwn/article-api/services/repository/mongodb"
	"github.com/zackwn/article-api/services/security"
	"github.com/zackwn/article-api/usecase"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Panic(err)
	}
}

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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password
		DB:       0,  // default DB
	})

	// services
	passwordHasher := security.NewPasswordHasher()
	authProvider := security.NewAuthProvider()
	fileStorage := filestorage.NewFileStorage()
	fph := fph.NewForgotPasswordHandler(redisClient)
	emailService := emailservice.NewEmailService()

	// repositories
	userRepository := mongodb.NewUserRepository(db)
	articleRepository := mongodb.NewArticleRepository(db)

	// usecases
	createUserUseCase := usecase.NewCreateUserUseCase(userRepository, passwordHasher)
	userLoginUseCase := usecase.NewUserLoginUseCase(userRepository, passwordHasher, authProvider)
	createArticleUseCase := usecase.NewCreateArticleUseCase(authProvider, articleRepository, userRepository)
	listArticlesUseCase := usecase.NewListArticlesUseCase(articleRepository, userRepository)
	forgotPasswordUseCase := usecase.NewForgotPasswordUseCase(userRepository, fph, emailService)

	// controllers
	userSignupController := c.NewUserSignupController(createUserUseCase, fileStorage)
	userSigninController := c.NewUserSigninController(userLoginUseCase)
	createArticleController := c.NewCreateArticleController(createArticleUseCase)
	listArticlesController := c.NewListArticlesController(listArticlesUseCase)
	forgotPasswordController := c.NewForgotPasswordController(forgotPasswordUseCase)

	http.HandleFunc("/user/signup", adaptController("POST", userSignupController))
	http.HandleFunc("/user/signin", adaptController("POST", userSigninController))

	http.HandleFunc("/articles/create", adaptController("POST", createArticleController))
	http.HandleFunc("/articles/list", adaptController("GET", listArticlesController))

	http.HandleFunc("/user/forgot-password", adaptController("POST", forgotPasswordController))

	http.Handle("/pictures/", http.StripPrefix("/pictures/", http.FileServer(http.Dir("./uploads"))))

	fmt.Println("ready")
	http.ListenAndServe(":8080", nil)
}
