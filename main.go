package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"github.com/gaisameer/auth-system/user"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func signUp (w http.ResponseWriter, r *http.Request) {
	var user user.User

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return
	}

	err = user.AddUser(collection)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func signIn (w http.ResponseWriter, r *http.Request) {
	var user user.User

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return
	}

	err = user.VerifyUser(collection)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func update (w http.ResponseWriter, r *http.Request) {
	type updater struct {
		user.User
		NewPassword string `json:"newPassword"`
	}

	var update updater

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &update)
	if err != nil {
		return
	}

	err = update.User.ChangePassword(collection, update.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

}

func delete (w http.ResponseWriter, r *http.Request) {
	var user user.User

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return
	}

	err = user.DeleteUser(collection)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

var client *mongo.Client
var collection *mongo.Collection

func main() {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongodb"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("user").Collection("credentials")
	defer client.Disconnect(ctx)

	http.HandleFunc("/signin", signIn)
	http.HandleFunc("/signup", signUp)
	http.HandleFunc("/update", update)
	http.HandleFunc("/delete", delete)
	log.Println(http.ListenAndServe(":8080", nil))
	
}
