package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	c "github.com/zackwn/article-api/controllers"
	filestorage "github.com/zackwn/article-api/services/file_storage"
	"github.com/zackwn/article-api/services/repository/mongodb"
	"github.com/zackwn/article-api/services/security"
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

	// services
	passwordHasher := security.NewPasswordHasher()
	authProvider := security.NewAuthProvider()
	fileStorage := filestorage.NewFileStorage()

	// repositories
	userRepository := mongodb.NewUserRepository(db)
	articleRepository := mongodb.NewArticleRepository(db)

	// usecases
	createUserUseCase := usecase.NewCreateUserUseCase(userRepository, passwordHasher)
	userLoginUseCase := usecase.NewUserLoginUseCase(userRepository, passwordHasher, authProvider)
	createArticleUseCase := usecase.NewCreateArticleUseCase(authProvider, articleRepository, userRepository)
	listArticles := usecase.NewListArticlesUseCase(articleRepository, userRepository)

	// controllers
	userSignupController := c.NewUserSignupController(createUserUseCase, fileStorage)
	userSigninController := c.NewUserSigninController(userLoginUseCase)
	createArticleController := c.NewCreateArticleController(createArticleUseCase)
	listArticlesController := c.NewListArticlesController(listArticles)

	http.HandleFunc("/user/signup", adaptController("POST", userSignupController))
	http.HandleFunc("/user/signin", adaptController("POST", userSigninController))
	http.HandleFunc("/articles/create", adaptController("POST", createArticleController))
	http.HandleFunc("/articles/list", adaptController("GET", listArticlesController))

	http.Handle("/pictures/", http.StripPrefix("/pictures/", http.FileServer(http.Dir("./uploads"))))

	http.HandleFunc("/test", func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		err := req.ParseMultipartForm(2 * 1024 * 1024)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("value name:", req.MultipartForm.Value["name"])
		f := req.MultipartForm.File["profile_picture"][0]
		fmt.Println(f.Filename)
		fmt.Println(f.Header["Content-Type"])
		file, err := f.Open()
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dst, err := os.Create(filepath.Join("./uploads/", f.Filename))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		bs, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = dst.Write(bs)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Println("ready")
	http.ListenAndServe(":8080", nil)
}
