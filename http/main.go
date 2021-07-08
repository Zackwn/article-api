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

func adaptController(controller c.Controller) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		res := controller.Handle(req)
		w.Header().Set("Content-Type", "application/json")
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
	fmt.Println("Successfully connected to the Database.")

	// services
	userRepository := mongodb.NewUserRepository(client.Database("articledb"))
	passwordHasher := security.NewPasswordHasher()
	authProvider := security.NewAuthProvider()

	// user signup
	createUserUseCase := usecase.NewCreateUserUseCase(userRepository, passwordHasher)
	userSignupController := c.NewUserSignupController(createUserUseCase)

	// user signin
	userLoginUseCase := usecase.NewUserLoginUseCase(userRepository, passwordHasher, authProvider)
	userSigninController := c.NewUserSigninController(userLoginUseCase)

	// test auth
	testAuthUseCase := usecase.NewTestAuthUseCase(authProvider)
	testAuthController := c.NewTestAuthController(testAuthUseCase)

	http.HandleFunc("/user/signup", adaptController(userSignupController))
	http.HandleFunc("/user/signin", adaptController(userSigninController))
	http.HandleFunc("/testauth", adaptController(testAuthController))

	fmt.Println("Server started...")
	http.ListenAndServe(":8080", nil)
}
