package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

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
	Username  string             `json:"userName,omitempty" bson:"userName,omitempty"`
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

//CreateJokeEndpoint is used for registration.
func CreateJokeEndpoint(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))
	response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	response.Header().Set("content-type", "application/json")
	var joke Joke
	_ = json.NewDecoder(request.Body).Decode(&joke)
	collection := client.Database("projectDb").Collection("jokes")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, _ := collection.InsertOne(ctx, joke)
	json.NewEncoder(response).Encode(result)
}

//DeleteJokeEndpoint is used for deleting joke.
func DeleteJokeEndpoint(response http.ResponseWriter, request *http.Request) {
	fmt.Println("delete")
	response.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	response.Header().Set("Access-Control-Allow-Methods", "OPTIONS,DELETE")
	response.Header().Set("content-type", "application/json")
	var joke Joke
	_ = json.NewDecoder(request.Body).Decode(&joke)
	collection := client.Database("projectDb").Collection("jokes")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, _ := collection.DeleteOne(ctx, Joke{ID: joke.ID})
	json.NewEncoder(response).Encode(result)
}

// GetAllJokesEndpoint returns all the jokes...
func GetAllJokesEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))
	response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	response.Header().Set("content-type", "application/json")
	var jokes []Joke
	collection := client.Database("projectDb").Collection("jokes")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var joke Joke
		cursor.Decode(&joke)
		jokes = append(jokes, joke)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(jokes)
}

//OptionsEndpoint is method .
func OptionsEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))
	response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	fmt.Println("options")
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
	router.HandleFunc("/joke", CreateJokeEndpoint).Methods("POST")
	router.HandleFunc("/joke", GetAllJokesEndpoint).Methods("GET")
	router.HandleFunc("/joke", DeleteJokeEndpoint).Methods("DELETE")
	router.HandleFunc("/joke", OptionsEndpoint).Methods("OPTIONS")
	http.ListenAndServe(":5000", router)
}
