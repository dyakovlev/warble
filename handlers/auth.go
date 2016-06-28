package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dyakovlev/warble/models"
	"github.com/dyakovlev/warble/utils"
)

func GetAuthHandler(c *gin.Context) {
	session, _ := c.MustGet("session").(models.Session)

	redir := c.DefaultQuery("redir", "/")

	// se url param expires the session (clears session cookie)

	if c.DefaultQuery("se", "") != "" {
		// TODO log
		utils.ExpireSessionCookie(c)
		c.HTML(http.StatusOK, "auth", gin.H{
			"login":    false,
			"register": true,
			"email":    nil,
		})
		return
	}

	// why would the user be logged in? handle this case anyway..

	if session.Auth != false {
		c.Redirect(http.StatusSeeOther, redir)
		return
	}

	// inactive sessions (user logged out or session expired) get a log-back-in screen

	if session.Auth == true && session.Uid != 0 {
		user := models.User{Res: session.Res}
		user.Load(session.Uid)
		c.HTML(http.StatusOK, "auth", gin.H{
			"login":    true,
			"register": false,
			"email":    utils.ObfuscateEmail(user.Email),
		})
		return
	}

	// sessions without an associated uid don't have an account associated with them,
	// so maybe they should make an account.

	c.HTML(http.StatusOK, "register", gin.H{
		"login":    true,
		"register": true,
		"email":    nil,
	})
}

func DoLogoutHandler(c *gin.Context) {
	session, _ := c.MustGet("session").(models.Session)

	// TODO log

	// TODO flash message

	session.Expire()

	c.Redirect(http.StatusSeeOther, "/")
}

func DoAuthHandler(c *gin.Context) {
	session, _ := c.MustGet("session").(models.Session)

	var err error
	user := models.User{Res: session.Res}

	if session.Uid != 0 {
		err = user.Load(session.Uid)
	} else {
		err = user.LoadByEmail(c.PostForm("email"))
	}

	if err != nil {
		// TODO log
		// TODO add no-user message to flash
		c.Redirect(http.StatusSeeOther, "/auth")
	}

	if user.Email != c.PostForm("email") {
		// TODO log
		// TODO add message
		session.Detach()
	}

	if utils.VerifyPass(user.Pass, c.PostForm("password")) {
		err = session.Authorize(user.Id)
		utils.SetSessionCookie(c, session.Res.Encid(session.Id))
		// TODO add logged-in message to flash
		// TODO log
		// TODO maybe makes sense to redirect to last-associated-project in session?
		c.Redirect(http.StatusSeeOther, c.DefaultQuery("redir", "/"))
	} else {
		// TODO log password auth failure
		c.Redirect(http.StatusOK, "/auth")
	}
}

func DoNewAccountHandler(c *gin.Context) {
	session, _ := c.MustGet("session").(models.Session)

	email := c.PostForm("email")
	pass := c.PostForm("password")

	// TODO validate email, pass

	user := models.User{
		Email: email,
		Pass:  utils.EncryptPass(pass),
		Res:   session.Res,
	}

	user.Store()

	session.Authorize(user.Id)

	utils.SetSessionCookie(c, session.Res.Encid(session.Id))

	c.Redirect(http.StatusSeeOther, "/profile")
}
