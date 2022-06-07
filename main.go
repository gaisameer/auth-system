package main

import (
	"net/http"
	"log"
	"encoding/json"
	"io/ioutil"
	"github.com/gaisameer/auth-system/user"
)

func signUp(w http.ResponseWriter, r *http.Request){
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

	err = user.AddUser()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func signIn(w http.ResponseWriter, r *http.Request) {
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

	err = user.VerifyUser()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func main() {

	http.HandleFunc("/signin", signIn)
	http.HandleFunc("/signup", signUp)
	log.Println(http.ListenAndServe(":8080", nil))
}