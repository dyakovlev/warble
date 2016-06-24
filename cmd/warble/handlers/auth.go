package main

import (
	"net/http"

	"github.com/gin-tonic/gin"
)

func GetAuthHandler(c *gin.Context) {
	session := c.Get("session")
	redir := c.DefaultQuery("redir", "/")

	// se url param expires the session (clears session cookie)

	if c.DefaultQuery("se", nil) != nil {
		// TODO log
		ExpireCookie(c, sessionCookie)
		c.HTML(http.StatusOk, "auth", gin.H{
			login:    false,
			register: true,
			email:    nil,
		})
		return
	}

	// why would the user be logged in? handle this case anyway..

	if session.auth {
		c.Redirect(http.StatusSeeOther, redir)
		return
	}

	// inactive sessions (user logged out or session expired) get a log-back-in screen

	if !session.auth && session.uid {
		user := User{resource: session.resource}
		user.load(session.uid)
		c.HTML(http.StatusOk, "auth", gin.H{
			login:    true,
			register: false,
			email:    ObfuscateEmail(user.email),
		})
		return
	}

	// sessions without an associated uid don't have an account associated with them,
	// so maybe they should make an account.

	c.HTML(http.StatusOk, "register", gin.H{
		login:    true,
		register: true,
		email:    nil,
	})
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

	// if there's an associated user in the session, load by it

	if session.uid != nil {
		err = user.load(session.uid)
	} else {
		err = user.loadByEmail(email)
	}

	if err == NoSuchUser {
		// TODO log
		// TODO add no-user message to flash
		c.Redirect(http.StatusSeeOther, "/auth")
	}

	if user.email != c.PostForm("email") {
		// TODO log
		// TODO add message
		session.detach()
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

	// TODO add success message to flash

	c.Redirect(http.StatusSeeOther, "/profile")
}
