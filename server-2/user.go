package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	Email    string `json:"email" bson:"email,omitempty"`
	Password string `json:"password" bson:"password,omitempty"`
	Salt     string `json:"salt" bson:"salt,omitempty"`
}

const usersCollectionsName = "users"

func (user *User) Insert() error {
	user.Password = GetMD5Hash(user.Password, user.Salt)

	userCollection := mongodb.Collection(usersCollectionsName)
	_, err := userCollection.InsertOne(context.TODO(), user)
	return err
}

func GetUserByEmail(email string) (user User, err error) {
	userCollections := mongodb.Collection(usersCollectionsName)
	filter := bson.D{{"email", email}}
	cursor := userCollections.FindOne(context.TODO(), filter)
	if err := cursor.Decode(&user); err != nil {
		return User{}, err
	}
	return
}
