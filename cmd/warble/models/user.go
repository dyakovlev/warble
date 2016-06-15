package main

import (
	"database/sql"
	"net/http"
)

// a User model represents a person that can be logged in to the site
type User struct {
}

func (m Session) name() string { return "user" }

func (m *User) Initialize(db *sql.DB, r *http.Request) *User {

}

func (m *User) GetById(id int) *User {

}

func (m *User) Upsert(r Row) (success bool) {

}
