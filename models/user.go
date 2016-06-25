package models

import (
	"github.com/gin-gonic/gin"
)

// a User model represents a person that can be logged in to the site
type User struct {
	// in schema
	Id    int
	Email string // used as login name
	Pass  string // salted, hashed password

	// not in schema
	Res *Resource // ref to initialized resources
}

func InitUser(res *Resource, c *gin.Context) (user *User, err error) {
	return &User{}, nil
}

func (u *User) Load(id int) (err error) {
	return nil
}

func (u *User) Store() error {
	return nil
}
