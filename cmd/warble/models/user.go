package main

import (
	"net/http"
)

// a User model represents a person that can be logged in to the site
type User struct {
	// in schema
	id    int
	email string
	salt  string // salt used to generate pass column
	pass  string // salted, hashed password

	// not in schema
	db *Database
}

func InitUser(db *Database, r *http.Request) (user *User, err error) {

}

func (u *User) load(id int) (err error) {

}

func (u *User) loadByName(username string) (err error) {

}

func (u *User) Store() (id int) {

}
