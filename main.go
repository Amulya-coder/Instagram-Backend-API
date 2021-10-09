package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Id       string
	Name     string
	Email    string
	Password string
}

type Post struct {
	Id               string
	Caption          string
	Image_URL        string
	Posted_Timestamp string
	UserId           string
}

func connectToMongo() (*mongo.Collection, *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	collection := client.Database("amulya").Collection("user")
	collection2 := client.Database("amulya").Collection("post")
	return collection, collection2
}

var collection1, collection2 *mongo.Collection = connectToMongo()

func userCreate(r *http.Request) *User {
	// Declare a new Person struct.
	var p User

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	// Do something with the Person struct...
	// fmt.Fprintf(w, "User: %+v", p)
	return &p
}

func postCreate(r *http.Request) *Post {
	// Declare a new Person struct.
	var p Post

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	// Do something with the Person struct...
	// fmt.Fprintf(w, "Post: %+v", p)
	return &p
}

func users(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		return
	}
	user := userCreate(req)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection1.InsertOne(ctx, bson.M{
		"Id":       user.Id,
		"Name":     user.Name,
		"Email":    user.Email,
		"Password": user.Password,
	})
	if err != nil {
		w.Write([]byte("Error"))
	}
	// id := res.InsertedID
	w.Write([]byte("Worked"))
}

func user(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		return
	}
	id := strings.TrimPrefix(req.URL.Path, "/user/")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := collection1.FindOne(ctx, bson.M{"Id": id})
	if doc == nil {
		w.Write([]byte("Couldn't find"))
	}
	res, _ := json.Marshal(doc)
	fmt.Fprintf(w, "Post: %+v", string(res))
}

func posts(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		return
	}
	post := postCreate(req)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection2.InsertOne(ctx, bson.M{
		"Id":               post.Id,
		"Caption":          post.Caption,
		"Image_URL":        post.Image_URL,
		"Posted_Timestamp": post.Posted_Timestamp,
		"UserId":           post.UserId,
	})
	if err != nil {
		w.Write([]byte("Error"))
	}
	// id := res.InsertedID
	w.Write([]byte("Worked"))
}

func post(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		return
	}
	id := strings.TrimPrefix(req.URL.Path, "/post/")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := collection2.FindOne(ctx, bson.M{"Id": id})

	if doc == nil {
		w.Write([]byte("Couldn't find"))
	}

	fmt.Fprintf(w, "Post: %+v", doc)

}

func postsUsers(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id := strings.TrimPrefix(req.URL.Path, "/posts/users/")
	options := options.Find()
	options.SetLimit(10)
	options.SetSkip(0)
	doc, _ := collection2.Find(ctx, bson.M{"UserId": id}, options)
	if doc == nil {
		w.Write([]byte("Couldn't find"))
	}

	fmt.Fprintf(w, "Posts: %+v", doc)
}

func main() {

	http.HandleFunc("/users", users)
	http.HandleFunc("/user/", user)
	http.HandleFunc("/posts", posts)
	http.HandleFunc("/post/", post)
	http.HandleFunc("/posts/users/", postsUsers)

	http.ListenAndServe(":8000", nil)
}
