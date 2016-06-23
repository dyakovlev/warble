package main

import (
	"github.com/gin-tonic/gin"
)

// get user profile page
func GetUserHandler(c *gin.Context) {
	session := c.Get("session")
	user := c.Get("user")

	// get list of all projects

	c.HTML(http.StatusOK, "user.tmpl.html", env)
}

// save profile info
func SaveUserHandler(c *gin.Context) {
	session := c.MustGet("session")
	user := c.MustGet("user")

}
