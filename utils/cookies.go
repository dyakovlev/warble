package utils

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
	SetCookie(c, name, "", time.Now())
}

func ExpireSessionCookie(c *gin.Context) {
	ExpireCookie(c, sessionCookie)
}

func SetCookie(c *gin.Context, name string, value string, expiration time.Time) {
	cookie := http.Cookie{
		Name:  name,
		Value: value,
		Path:  "/"}

	http.SetCookie(c.Writer, &cookie)
}

func GetSessionCookie(c *gin.Context) (string, error) {
	return c.Cookie(sessionCookie)
}

func SetSessionCookie(c *gin.Context, encSid string) {
	exp := time.Now().AddDate(0, 1, 0) // 1 month (TODO extend session length automatically?)
	SetCookie(c, sessionCookie, encSid, exp)
}
