package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-tonic/gin"
)

// environment variables
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

	resource := Resource{
		db:      sql.Open("postgres", os.Getenv(postgresURL)),
		crypter: NewIDCodec(os.Getenv(encIDKey)),
	}

	router.GET("/about", staticPage("about"))

	router.GET("/status", resource.withSession(), InGroup(Admin), StatusHandler)

	router.GET("/", staticPage("index"))

	router.GET("/auth", resource.withSession(), GetAuthHandler)
	router.POST("/auth/new", resource.withSession(), DoNewAccountHandler)
	router.POST("/auth/login", rsource.withSession(), DoAuthHandler)
	router.POST("/auth/logout", resource.withSession(), DoLogoutHandler)

	// logged-in app endpoints
	app := router.Group("/", resource.withSession(), InGroup(User))
	{
		// modify user profile
		router.GET("/user", resource.withUser(), GetUserHandler)
		router.POST("/user", resource.withUser(), SaveUserHandler)

		// modify project
		router.GET("/project", resource.withProject(), GetProjectHandler)
		router.POST("/project", resource.withProject(), SaveProjectHandler)

		// modify clip audio file
		router.GET("/clip", resource.withClip(), GetClipHandler)
		router.POST("/clip", resource.withClip(), SaveClipHandler)
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
