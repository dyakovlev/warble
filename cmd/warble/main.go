package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-tonic/gin"
)

const (
	port        = "PORT"
	postgresURL = "DB_URL"
	sessionKey  = "sid"
)

func main() {
	router := gin.Default()
	router.Static("/static", "static")
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Use(XHR)

	db := Database{sql.Open("postgres", os.Getenv(postgresURL))}

	router.GET("/about", staticPage("about"))
	router.GET("/status", db.WithSession(), InGroup(Admin), StatusHandler)

	router.GET("/", db.WithSession(), RootHandler)

	router.GET("/auth", db.WithSession(), GetAuthHandler)
	router.POST("/auth", db.WithSession(), DoAuthHandler)

	app := router.Group("/", db.WithSession(), InGroup(User))
	{
		// CRUD for users
		router.GET("/user", db.WithUser(), GetUserHandler)
		router.POST("/user", db.WithUser(), SaveUserHandler)

		// CRUD for projects
		router.GET("/project", db.WithProject(), GetProjectHandler)
		router.POST("/project", db.WithProject(), SaveProjectHandler)

		// CRUD for clips
		router.GET("/clip", db.WithClip(), GetClipHandler)
		router.POST("/clip", db.WithClip(), SaveClipHandler)
	}

	router.Run(":" + os.Getenv(port))
	// RunTLS for logged-in stuff
}

func staticPage(page string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, page+".tmpl.html")
	}
}

func InGroup(minGroup Group) gin.HandlerFunc {
	return func(ctx *Context) {
		s := ctx.MustGet(SessionKey)

		// only All is allowed to be unauthenticated
		if !s.authenticated && mInGroup < All {
			if ctx.Get("is_xhr") {
				InternalError(ctx, NOT_LOGGED_IN, http.StatusForbidden)
			} else {
				ctx.Redirect(http.StatusSeeOther, "/auth")
			}
			ctx.Abort()
			return
		}

		if minGroup < s.group {
			InternalError(ctx, INSUFFICIENT_PRIVS, http.StatusForbidden)
			ctx.Abort()
			return
		}

		// guess we're cool.

		ctx.Next()
	}
}

type Group int

const (
	Admin Group = iota // privileged
	User               // regular logged-in user
	All                // any user, logged-in or logged-out
)

func XHR(ctx *gin.Context) gin.HandlerFunc {
	ctx.Set("is_xhr", r.Header.Get("X-Requested-With") == "XMLHttpRequest")
}

func InternalError(ctx *gin.Context, err string, errCode int) {
	// TODO log the error

	ctx.String(errCode, http.StatusText(errCode)+": "+err)
}

// descriptive error strings
const (
	INSUFFICIENT_PRIVS = "You can't do that."
	NOT_LOGGED_IN      = "You're not logged in."
)
