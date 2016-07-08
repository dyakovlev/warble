package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func ExpireCookie(c *gin.Context, name string) {
	SetCookie(c, name, "", time.Now())
}

func SetCookie(c *gin.Context, name string, value string, expiration time.Time) {
	cookie := http.Cookie{Name: name, Value: value, Path: "/"}
	http.SetCookie(c.Writer, &cookie)
}

func GetCookie(c *gin.Context, name string) (string, error) {
	cookie, err = c.Cookie(name)

	if err != nil {
		Error("GetCookie: failed to retrieve cookie", name, ":", err)
	}

	return cookie, err
}

const sessionCookie = "s"

func GetSessionCookie(c *gin.Context) (string, error) {
	return GetCookie(sessionCookie)
}

func SetSessionCookie(c *gin.Context, encSid string) {
	exp := time.Now().AddDate(0, 1, 0) // 1 month (TODO extend session length automatically?)
	SetCookie(c, sessionCookie, encSid, exp)
}

func ExpireSessionCookie(c *gin.Context) {
	ExpireCookie(c, sessionCookie)
}
