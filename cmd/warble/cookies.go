package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// cookie names
const (
	sessionCookie = "s"
)

func ExpireCookie(c *gin.Context, name string) {
	SetCookie(c, name, "", time.Now(), false)
}

func SetCookie(c *gin.Context, name string, value string, expiration time.Time, secure bool) {
	cookie := http.Cookie{
		Name:       name,
		Value:      value,
		Path:       "/",
		Expires:    expiration,
		RawExpires: expiration.Format(time.UnixDate),
		Secure:     secure,
		HttpOnly:   false,
	}
	http.SetCookie(c.Writer, &cookie)
}

func SetSessionCookie(c *gin.Context, s *Session) {
	encSid := s.res.crypter.encid(s.id)
	exp := time.Now().AddDate(0, 1, 0) // 1 month (TODO extend session length automatically?)
	SetCookie(c, sessionCookie, encSid, exp, true)
}
