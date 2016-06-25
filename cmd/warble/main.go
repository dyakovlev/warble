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

	router.GET("/", staticPage("index"))

	router.GET("/about", staticPage("about"))

	router.GET("/status", SessionMiddleware(r), InGroup(Admin), handlers.StatusHandler)

	auth := router.Group("/auth", SessionMiddleware(r))
	{
		router.GET("/auth", handlers.GetAuthHandler)
		router.POST("/auth/new", handlers.DoNewAccountHandler)
		router.POST("/auth/login", handlers.DoAuthHandler)
		router.POST("/auth/logout", handlers.DoLogoutHandler)
	}

	app := router.Group("/", SessionMiddleware(r), InGroup(User))
	{
		router.GET("/user", handlers.GetUserHandler)
		router.POST("/user", handlers.SaveUserHandler)

		router.GET("/project", handlers.GetProjectHandler)
		router.POST("/project", handlers.SaveProjectHandler)

		router.GET("/clip", handlers.GetClipHandler)
		router.POST("/clip", handlers.SaveClipHandler)
	}

	router.Run(":" + os.Getenv(Port))
	// TODO RunTLS for logged-in stuff
}

func XHRMiddleware(ctx *gin.Context) {
	ctx.Set("is_xhr", r.Header.Get("X-Requested-With") == "XMLHttpRequest")
}

func SessionMiddleware(r *models.Resource) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s, err := models.InitSession(r, ctx)

		// TODO handle err

		ctx.Set("session", s)
	}
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
