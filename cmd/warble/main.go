package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/dyakovlev/warble/handlers"
	"github.com/dyakovlev/warble/utils"
	"github.com/gin-gonic/gin"
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

	resource := Resource{}
	var err error

	resource.db, err = sql.Open("postgres", os.Getenv(PostgresURL))
	resource.crypter = utils.NewIDCodec(os.Getenv(EncIDKey))

	router.GET("/about", staticPage("about"))

	router.GET("/status", resource.withSession(), InGroup(Admin), handlers.StatusHandler)

	router.GET("/", staticPage("index"))

	router.GET("/auth", resource.withSession(), handlers.GetAuthHandler)
	router.POST("/auth/new", resource.withSession(), handlers.DoNewAccountHandler)
	router.POST("/auth/login", resource.withSession(), handlers.DoAuthHandler)
	router.POST("/auth/logout", resource.withSession(), handlers.DoLogoutHandler)

	// logged-in app endpoints
	app := router.Group("/", resource.withSession(), InGroup(User))
	{
		// modify user profile
		router.GET("/user", resource.withUser(), handlers.GetUserHandler)
		router.POST("/user", resource.withUser(), handlers.SaveUserHandler)

		// modify project
		router.GET("/project", resource.withProject(), handlers.GetProjectHandler)
		router.POST("/project", resource.withProject(), handlers.SaveProjectHandler)

		// modify clip audio file
		router.GET("/clip", resource.withClip(), handlers.GetClipHandler)
		router.POST("/clip", resource.withClip(), handlers.SaveClipHandler)
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
