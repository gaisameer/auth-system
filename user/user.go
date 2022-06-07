package user

import (
	"log"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var dict = make(map[string]string)

func (user User)checkIfPresent() bool {
	_, ok := dict[user.Username]
	return ok
}

func (user User)AddUser() error {
	if user.checkIfPresent() {
		return errors.New("User already exist")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		return err
	}
	dict[user.Username] = string(hash)
	log.Printf("User %s added\n", user.Username)
	return nil
}

func (user User)VerifyUser() error {
	if user.checkIfPresent() == false {
		return errors.New("User doesn't exist")
	}

	err := bcrypt.CompareHashAndPassword([]byte(dict[user.Username]), []byte(user.Password))
	if err != nil {
		return errors.New("Wrong credentials")
	}
}