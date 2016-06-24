package main

import (
	"net/http"

	"github.com/gin-tonic/gin"
)

func GetAuthHandler(c *gin.Context) {
	session := c.Get("session")
	redir := c.DefaultQuery("redir", "/")

	// why would the user be logged in? handle this case anyway..

	if session.auth {
		c.Redirect(http.StatusSeeOther, redir)
		return
	}

	// inactive sessions (user logged out or session expired) get a log-back-in screen
	// with their email prefilled

	if !session.auth && session.uid {
		c.HTML(http.StatusOk, "login", userParams)
	}

	// sessions without an associated uid don't have an account associated with them..
	// maybe they should make an account.

	c.HTML(http.StatusOk, "new_account")
}

func DoLogoutHandler(c *gin.Context) {
	session := c.Get("session")

	// TODO log

	session.expire()

	c.Redirect(http.StatusSeeOther, "/")
}

func DoAuthHandler(c *gin.Context) {
	session := c.Get("session")

	var err error

	user := User{resource: session.resource}

	err = user.loadByName(session.uid, c.PostForm("username"))

	switch {
	case err == AuthUsernameSessionMismatch:
		// TODO log
		session.detach()
	case err == NoSuchUser:
		// TODO log
		// TODO add no-user message to flash
		c.Redirect(http.StatusSeeOther, "/auth")
	}

	if verifyPass(user.pass, c.PostForm("password")) {
		err = session.authorize(user.id)
		// TODO add logged-in message to flash
		// TODO log
		// TODO maybe makes sense to redirect to last-associated-project in session?
		c.Redirect(http.StatusSeeOther, c.DefaultQuery("redir", "/"))
	} else {
		// TODO log password auth failure
		c.Redirect(http.StatusOk, "/auth")
	}
}

func DoNewAccountHandler(c *gin.Context) {
	session := c.Get("session")

	email := c.PostForm("email")
	pass := c.PostForm("password")

	// TODO validate email, pass

	user := User{
		email:    email,
		pass:     encryptPass(pass),
		resource: session.res,
	}

	user.store()

	session.authorize(user.id)

	// TODO add success message

	c.Redirect(http.StatusSeeOther, "/profile")
}
