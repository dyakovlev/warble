package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/dyakovlev/warble/handlers"
	"github.com/dyakovlev/warble/models"
	"github.com/dyakovlev/warble/utils"
)

// environment variables
const (
	Port        = "PORT"
	PostgresURL = "DATABASE_URL"
	EncIDKey    = "ENCID_KEY"
)

func main() {
	router := gin.Default()
	router.Static("/static", "static")
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Use(XHRMiddleware)

	r, err := models.NewResource(os.Getenv(PostgresURL), os.Getenv(EncIDKey))

	router.GET("/about", staticPage("about"))

	router.GET("/status", r.WithSession(), InGroup(Admin), handlers.StatusHandler)

	router.GET("/", staticPage("index"))

	router.GET("/auth", r.WithSession(), handlers.GetAuthHandler)
	router.POST("/auth/new", r.WithSession(), handlers.DoNewAccountHandler)
	router.POST("/auth/login", r.WithSession(), handlers.DoAuthHandler)
	router.POST("/auth/logout", r.WithSession(), handlers.DoLogoutHandler)

	// logged-in app endpoints
	app := router.Group("/", r.WithSession(), InGroup(User))
	{
		// modify user profile
		router.GET("/user", r.WithUser(), handlers.GetUserHandler)
		router.POST("/user", r.WithUser(), handlers.SaveUserHandler)

		// modify project
		router.GET("/project", r.WithProject(), handlers.GetProjectHandler)
		router.POST("/project", r.WithProject(), handlers.SaveProjectHandler)

		// modify clip audio file
		router.GET("/clip", r.WithClip(), handlers.GetClipHandler)
		router.POST("/clip", r.WithClip(), handlers.SaveClipHandler)
	}

	router.Run(":" + os.Getenv(Port))
	// TODO RunTLS for logged-in stuff
}

func XHRMiddleware(ctx *gin.Context) {
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
		session, ok := s.(models.Session)

		// only All is allowed to be unauthenticated
		if !s.Auth && minGroup < All {
			if ctx.Get("is_xhr") {
				InternalError(ctx, NOT_LOGGED_IN, http.StatusForbidden)
			} else {
				// TODO last page url param to redir to cur page after auth
				ctx.Redirect(http.StatusSeeOther, "/auth")
			}
			ctx.Abort()
			return
		}

		if minGroup < s.Group {
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
