package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dyakovlev/warble/models"
	"github.com/dyakovlev/warble/utils"
)

func GetAuthHandler(c *gin.Context) {
	session, _ := c.MustGet("session").(*models.Session)

	redir := c.DefaultQuery("redir", "/")

	// se url param expires the session (clears session cookie)

	if c.DefaultQuery("se", "") != "" {
		// TODO log
		utils.ExpireSessionCookie(c)
		c.HTML(http.StatusOK, "auth.tmpl.html", gin.H{
			"login":    false,
			"register": true,
			"email":    nil,
		})
		return
	}

	// why would the user be logged in? handle this case anyway..

	if session.Auth == true {
		c.Redirect(http.StatusSeeOther, redir)
		return
	}

	// inactive sessions (user logged out or session expired) get a log-back-in screen

	if session.Auth == false && session.Uid != 0 {
		if user, err := models.InitUser(session); err == nil {
			c.HTML(http.StatusOK, "auth.tmpl.html", gin.H{
				"login":    true,
				"register": false,
				"email":    utils.ObfuscateEmail(user.Email),
			})
			return
		}
		// fall through to login/register if failed to load user..
	}

	// sessions without an associated uid don't have an account associated with them,
	// so maybe they should make an account.

	c.HTML(http.StatusOK, "auth.tmpl.html", gin.H{
		"login":    true,
		"register": true,
		"email":    nil,
	})
}

func DoLogoutHandler(c *gin.Context) {
	session, _ := c.MustGet("session").(*models.Session)

	// TODO log

	// TODO flash message

	session.Expire()

	c.Redirect(http.StatusSeeOther, "/")
}

func DoAuthHandler(c *gin.Context) {
	session, _ := c.MustGet("session").(*models.Session)

	var err error
	user := models.User{Res: session.Res}

	if err = user.LoadByEmail(c.PostForm("email")); err != nil {
		utils.Error("DoAuthHandler: failed to load user by email")
		// TODO add no-user message to flash
		c.Redirect(http.StatusSeeOther, "/auth")
		return
	}

	if utils.VerifyPass(user.Pass, c.PostForm("password")) {
		err = session.Authorize(&user)
		utils.SetSessionCookie(c, session.Res.Encid(session.Id))
		// TODO add logged-in message to flash
		c.Redirect(http.StatusSeeOther, c.DefaultQuery("redir", "/"))
	} else {
		utils.Error("DoAuthHandler: bad password")
		c.Redirect(http.StatusSeeOther, "/auth")
	}
}

func DoNewAccountHandler(c *gin.Context) {
	session, _ := c.MustGet("session").(*models.Session)

	email := c.PostForm("email")
	pass := c.PostForm("password")

	// TODO validate email, pass

	user := models.User{
		Email: email,
		Pass:  utils.EncryptPass(pass),
		Admin: false,
		Res:   session.Res,
	}

	if err := user.Store(); err != nil {
		utils.Error("DoNewAccountHandler error storing user: ", err)
		// TODO set error flash
		c.Redirect(http.StatusSeeOther, "/auth")
		return
	}

	session.Authorize(&user)
	utils.SetSessionCookie(c, session.Res.Encid(session.Id))
	c.Redirect(http.StatusSeeOther, "/profile")
}
