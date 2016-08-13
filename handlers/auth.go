package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dyakovlev/warble/models"
	"github.com/dyakovlev/warble/utils"
)

func ServeAuthPage(c *gin.Context) {
	session, _ := c.MustGet("session").(*models.Session)

	// logged in already

	if session.Auth == true {
		c.Redirect(http.StatusSeeOther, c.DefaultQuery("redir", "/"))
		return
	}

	// "se" query param, i.e. ignore session

	if c.DefaultQuery("se", "") != "" {
		utils.ExpireSessionCookie(c)
		c.HTML(http.StatusOK, "auth.tmpl.html", gin.H{
			"login":    true,
			"register": false,
			"email":    nil,
		})
		return
	}

	// existing session

	if session.Uid != 0 {
		if user, err := models.InitUser(session); err == nil {
			c.HTML(http.StatusOK, "auth.tmpl.html", gin.H{
				"login":    true,
				"register": false,
				"email":    utils.ObfuscateEmail(user.Email),
			})
			return
		}
	}

	// no session

	c.HTML(http.StatusOK, "auth.tmpl.html", gin.H{
		"login":    false,
		"register": true,
		"email":    nil,
	})
}

func DoLogout(c *gin.Context) {
	session, _ := c.MustGet("session").(*models.Session)
	session.Expire()
	c.Redirect(http.StatusSeeOther, "/")
}

func DoAuth(c *gin.Context) {
	session, _ := c.MustGet("session").(*models.Session)

	email, emailErr := utils.ParseEmail(c.PostForm("email"))
	if emailErr != nil {
		retryAuth(c, fmt.Sprintf("failed to parse email %v", email), "bad email supplied")
		return
	}

	pass, passErr := utils.ParsePassword(c.PostForm("password"))
	if passErr != nil {
		retryAuth(c, fmt.Sprintf("failed to parse password %v", pass), "bad password supplied")
		return
	}

	user := models.User{Res: session.Res}

	switch {
	case email == "" && session.Uid != 0:
		if err := user.Load(session.Uid); err != nil {
			retryAuth(c, fmt.Sprintf("failed to load user %v", session.Uid), "bad session user")
			return
		}
	case email != "" && session.Uid == 0:
		if err := user.LoadByEmail(email); err != nil {
			retryAuth(c, fmt.Sprintf("failed to load user by email %v", email), "bad email")
			return
		}
	default:
		retryAuth(c, "supplied both or neither session/email", "bad data")
		return
	}

	if !utils.VerifyPass(user.Pass, pass) {
		retryAuth(c, "bad password", "bad password")
		return
	}

	if err := session.Authorize(&user); err != nil {
		utils.SetSessionCookie(c, session.Res.Encid(session.Id))
	}

	c.Redirect(http.StatusSeeOther, c.DefaultQuery("redir", "/"))
}

func retryAuth(c *gin.Context, errLog string, errFlash string) {
	utils.Error("[DoAuthHandler]", errLog)
	// TODO set flash to errFlash
	c.Redirect(http.StatusSeeOther, "/auth")
}

func DoNewAccount(c *gin.Context) {
	session, _ := c.MustGet("session").(*models.Session)

	email, emailErr := utils.ParseEmail(c.PostForm("email"))
	if emailErr != nil {
		retryNewAccount(c, fmt.Sprintf("failed to parse email %v", email), "bad email supplied")
		return
	}

	pass, passErr := utils.ParsePassword(c.PostForm("password"))
	if passErr != nil {
		retryNewAccount(c, fmt.Sprintf("failed to parse password %v", pass), "bad password supplied")
		return
	}

	user := models.User{
		Email: email,
		Pass:  utils.EncryptPass(pass),
		Admin: false,
		Res:   session.Res,
	}

	if err := user.Store(); err != nil {
		retryNewAccount(c, fmt.Sprintf("didn't store user: %v", err), "didn't store user")
		return
	}

	if err := session.Authorize(&user); err != nil {
		utils.Error("failed to authorize user into session", err)
	}

	utils.SetSessionCookie(c, session.Res.Encid(session.Id))

	c.Redirect(http.StatusSeeOther, "/profile")
}

func retryNewAccount(c *gin.Context, errLog string, errFlash string) {
	utils.Error("[DoNewAccountHandler]", errLog)
	// TODO set flash to errFlash
	c.Redirect(http.StatusSeeOther, "/auth")
}
