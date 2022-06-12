package user

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// var dict = make(map[string]string)

func (user User) checkIfPresent(collection *mongo.Collection) (User, error) {
	
	var result User
	err := collection.FindOne(context.Background(), bson.D{{"username", user.Username}}).Decode(&result)
	return result, err
}

func (user User) AddUser(collection *mongo.Collection) error {
	_, err := user.checkIfPresent(collection)
	if err == nil {
		// if data couldnt be fetched, mongo.ErrNoDocument will be raised
		return errors.New("User already exist")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}

	log.Printf("User %s added with id %s\n", user.Username, result.InsertedID)
	return nil
}

func (user User) VerifyUser(collection *mongo.Collection) error {

	entry, err := user.checkIfPresent(collection)
	if err == mongo.ErrNoDocuments {
		return errors.New("User doesn't exist")
	}

	err = bcrypt.CompareHashAndPassword([]byte(entry.Password), []byte(user.Password))
	if err != nil {
		return errors.New("Wrong credentials")
	}

	log.Printf("User %s signed in", user.Username)
	return nil
}

func (user User) ChangePassword(collection *mongo.Collection, newPassword string) error {

	var newUser User
	err := user.VerifyUser(collection)
	if err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 8)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	err = collection.FindOneAndReplace(context.Background(), bson.D{{"username", user.Username}}, user).Decode(&newUser)
	if err == nil {
		log.Printf("Updated password of %s", newUser.Username)
	}
	return err
}

func (user User) DeleteUser(collection *mongo.Collection) error {

	var deletedUser User
	err := user.VerifyUser(collection)
	if err != nil {
		return err
	}

	err = collection.FindOneAndDelete(context.Background(), bson.D{{"username", user.Username}}).Decode(&deletedUser)
	if err == nil {
		log.Printf("Deleted user %s", deletedUser.Username)
	}
	return err
}
