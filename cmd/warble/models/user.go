package main

import (
	"net/http"
)

// a User model represents a person that can be logged in to the site
type User struct {
	// in schema
	id    int
	email string // used as login name
	pass  string // salted, hashed password

	// not in schema
	resource *Resource // ref to initialized resources
}

func InitUser(res *Resource, req *http.Request) (*User, error) {

}

func (u *User) load(id int) (err error) {

}

func (u *User) loadByName(uidFromSession int, username string) (err error) {

}

func (u *User) store() error {

}
