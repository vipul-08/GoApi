package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

//RandomResponse is a model for any random response.
type RandomResponse struct {
	resp   string
	status int
}

//User is a model for users.
type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Fullname string             `json:"fullName,omitempty" bson:"fullName,omitempty"`
	Username string             `json:"userName,omitempty" bson:"userName,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

//Joke is a model for jokes
type Joke struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Joketitle string             `json:"jokeTitle,omitempty" bson:"jokeTitle,omitempty"`
	Jokevalue string             `json:"jokeString,omitempty" bson:"jokeString,omitempty"`
	Userid    string             `json:"userId,omitempty" bson:"userId,omitempty"`
}

//CreateUserEndpoint is used for registration.
func CreateUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var user User
	_ = json.NewDecoder(request.Body).Decode(&user)
	collection := client.Database("projectDb").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(response).Encode(result)
}

//LoginUserEndpoint is used for login.
func LoginUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var user User
	_ = json.NewDecoder(request.Body).Decode(&user)

	var foundUser User
	collection := client.Database("projectDb").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, User{Username: user.Username}).Decode(&foundUser)
	if err != nil {
		response.Write([]byte(`{ "message": "No User Found!!" }`))
	} else {
		if foundUser.Password == user.Password {
			json.NewEncoder(response).Encode(foundUser)
		} else {
			response.Write([]byte(`{ "message": "Wrong Password!!" }`))
		}
	}
}

func main() {
	fmt.Println("Starting the Application....")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()

	router.HandleFunc("/user", CreateUserEndpoint).Methods("POST")
	router.HandleFunc("/login", LoginUserEndpoint).Methods("POST")
	http.ListenAndServe(":5000", router)
}
