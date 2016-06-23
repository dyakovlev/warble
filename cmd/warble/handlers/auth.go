package main

import (
	"github.com/gin-tonic/gin"
)

// validate auth info, write into session
func DoAuthHandler(c *gin.Context) {
	session := c.Get("session")
	user := User{resource: session.resource}

	// TODO if username is defined in session and doesn't match the one provided in request,
	// something is wrong and/or sketchy, should log and error

	if session.uid == nil {
		user.loadByName(name)
	} else {
		user.load(session.uid)
	}

	if verifyPass(user.pass, session.resource.crypter.encryptPass(user.salt, suppliedPass)) {
		err := session.authorize(user.id)
		// redir to previous page or /
	} else {
		// log
		// retry
	}
}

// rubber stamp a new user and project
func DoNewAccountHandler(c *gin.Context) {
	session := c.Get("session")

}
