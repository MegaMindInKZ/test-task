package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

var ctx context.Context
var mongodb *mongo.Database

const mongoURL = "mongodb://localhost:27017"
const dbName = "task"

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer client.Disconnect(ctx)

	mongodb = client.Database(dbName)

	command := bson.D{{"create", usersCollectionsName}}
	var result bson.M
	if err := mongodb.RunCommand(context.TODO(), command).Decode(&result); err != nil {
		fmt.Println(err)
	}

	r := chi.NewRouter()
	r.Post("/create-user", createUser)
	r.Get("/get-user/{email}", getUser)

	server := http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8002",
	}

	server.ListenAndServe()

	log.Printf("serving on port %d", 8002)

}

func createUser(writer http.ResponseWriter, request *http.Request) {
	var user User
	if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
		SendErrorMessage(writer, http.StatusInternalServerError, err, errMessageInternalServerError)
		return
	}

	v := NewValidator()
	v.Check(user.Email)
	if !v.isValid() {
		SendValidatorMessage(writer, v)
		return
	}

	salt, err := attachSalt()
	if err != nil {
		SendErrorMessage(writer, http.StatusInternalServerError, err, errMessageInternalServerError)
		return
	}

	user.Salt = salt

	if err := user.Insert(); err != nil {
		if err == ErrDuplicationEmail {
			SendErrorMessage(writer, http.StatusBadRequest, err, errMessageDuplicationEmail)
			return
		}
		SendErrorMessage(writer, http.StatusInternalServerError, err, errMessageInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusCreated)
}

func getUser(writer http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	user, err := GetUserByEmail(email)
	if err != nil {
		if err.Error() == ErrNoDocument.Error() {
			SendErrorMessage(writer, http.StatusBadRequest, err, errMessageNoDocument)
			return
		}
		SendErrorMessage(writer, http.StatusInternalServerError, err, errMessageInternalServerError)
		return
	}
	jsonResp, err := json.Marshal(user)
	if err != nil {
		SendErrorMessage(writer, http.StatusInternalServerError, err, errMessageInternalServerError)
		return
	}
	writer.Write(jsonResp)
}
