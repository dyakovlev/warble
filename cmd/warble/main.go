package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-tonic/gin"
)

// environment variable names
const (
	port        = "PORT"
	postgresURL = "DB_URL"
	encIDKey    = "ENCID_KEY"
)

func main() {
	router := gin.Default()
	router.Static("/static", "static")
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Use(XHRMiddleware)

	db := Database{sql.Open("postgres", os.Getenv(postgresURL)), NewIDCodec(os.Getenv(encIDKey))}

	router.GET("/about", staticPage("about"))

	router.GET("/status", db.withSession(), InGroup(Admin), StatusHandler)

	// GET root page, redir to active project if one is in session
	router.GET("/", db.withSession(), staticPage("index"))

	// GET login page, POST auth or new-account endpoints
	router.GET("/auth", db.withSession(), staticPage("auth"))
	router.POST("/auth", db.withSession(), DoAuthHandler)
	router.POST("/auth/new", db.withSession(), DoNewAccountHandler)

	// logged-in app endpoints
	app := router.Group("/", db.withSession(), InGroup(User))
	{
		// modify user profile
		router.GET("/user", db.withUser(), GetUserHandler)
		router.POST("/user", db.withUser(), SaveUserHandler)

		// modify project
		router.GET("/project", db.withProject(), GetProjectHandler)
		router.POST("/project", db.withProject(), SaveProjectHandler)

		// modify clip audio file
		router.GET("/clip", db.withClip(), GetClipHandler)
		router.POST("/clip", db.withClip(), SaveClipHandler)
	}

	router.Run(":" + os.Getenv(port))
	// TODO RunTLS for logged-in stuff
}

func XHRMiddleware(ctx *gin.Context) gin.HandlerFunc {
	ctx.Set("is_xhr", r.Header.Get("X-Requested-With") == "XMLHttpRequest")
}

func staticPage(page string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, page+".tmpl.html")
	}
}

func InGroup(minGroup Group) gin.HandlerFunc {
	return func(ctx *Context) {
		s := ctx.MustGet("session")

		// only All is allowed to be unauthenticated
		if !s.auth && minGroup < All {
			if ctx.Get("is_xhr") {
				InternalError(ctx, NOT_LOGGED_IN, http.StatusForbidden)
			} else {
				// TODO last page url param to redir to cur page after auth
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

		ctx.Next()
	}
}

type Group int

const (
	Admin Group = iota // privileged
	User               // regular logged-in user
	All                // any user, logged-in or logged-out
)
