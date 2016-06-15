package main

import (
	"net/http"
)

// a User model represents a person that can be logged in to the site
type User struct {
}

func InitUser(db *Database, r *http.Request) (user *User, err error) {

}

func LoadUser(id int) (user *User, err error) {

}

func (m *User) Store() (id int) {

}
