package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"posty/api/handler"
	"posty/pkg/post"
	"posty/pkg/user"

	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := DatabaseConnection()
	if err != nil {
		log.Fatal("Database Connection Error $s", err)
	}
	fmt.Println("DB connected")

	userCollection := db.Collection("user")
	userRepo := user.NewRepo(userCollection)
	userService := user.NewService(userRepo)

	postCollection := db.Collection("post")
	postRepo := post.NewRepo(postCollection)
	postService := post.NewService(postRepo)

	http.HandleFunc("/ping", handler.PingHandler)
	http.HandleFunc("/users", handler.AddUser(userService))
	http.HandleFunc("/users/", handler.GetUser(userService))

	http.HandleFunc("/posts", handler.AddPost(postService, userService))
	http.HandleFunc("/posts/", handler.GetPost(postService))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func DatabaseConnection() (*mongo.Database, error) {
	// log.Println(os.Getenv("DB_URL"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DB_URL")))
	if err != nil {
		return nil, err
	}
	db := client.Database("posty")
	return db, nil
}
