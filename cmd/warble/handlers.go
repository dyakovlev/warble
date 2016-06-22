package main

import (
	"github.com/gin-tonic/gin"
)

// dump stats
func StatusHandler(c *gin.Context) {
	// total sessions
	// active sessions
	// total users
	// total projects
	// total clips
	// requests seen today
	// memory stats
	// s3 stats
}

// validate auth info, write into session
func DoAuthHandler(c *gin.Context) {
	session := c.Get("session")
	user := User{db: session.db}

	// TODO if username is defined in session and doesn't match the one provided in request,
	// something is wrong and/or sketchy, should log and error

	if session.uid == nil {
		user.loadByName(name)
	} else {
		user.load(session.uid)
	}

	if verifyPass(user.pass, session.db.crypter.encryptPass(user.salt, suppliedPass)) {
		err := session.authorize(user.id)
		// redir to previous page or /
	} else {
		// log
		// retry
	}
}

// rubber stamp a new user and project
func DoNewAccountHandler(c *gin.Context) {

}

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

// get project page
func GetProjectHandler(c *gin.Context) {
	session := c.MustGet("session")
	project := c.MustGet("project")
}

// save project info
func SaveProjectHandler(c *gin.Context) {
	session := c.MustGet("session")
	project := c.MustGet("project")

}

// retrieve clip info
func GetClipHandler(c *gin.Context) {
	session := c.MustGet("session")
	clip := c.MustGet("clip")
}

// save clip info
func SaveClipHandler(c *gin.Context) {
	session := c.MustGet("session")
	clip := c.MustGet("clip")

}
