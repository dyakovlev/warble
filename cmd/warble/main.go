package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/dyakovlev/warble/handlers"
	"github.com/dyakovlev/warble/models"
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

	r, _ := models.NewResource(os.Getenv(PostgresURL), os.Getenv(EncIDKey))

	router.GET("/", staticPage("index"))

	router.GET("/about", staticPage("about"))

	router.GET("/status", SessionMiddleware(r), InGroup(Admin), handlers.StatusHandler)

	auth := router.Group("/auth", SessionMiddleware(r))
	{
		auth.GET("", handlers.GetAuthHandler)
		auth.POST("/new", handlers.DoNewAccountHandler)
		auth.POST("/login", handlers.DoAuthHandler)
		auth.POST("/logout", handlers.DoLogoutHandler)
	}

	app := router.Group("/", SessionMiddleware(r), InGroup(User))
	{
		app.GET("/user", handlers.GetUserHandler)
		app.POST("/user", handlers.SaveUserHandler)

		app.GET("/project", handlers.GetProjectHandler)
		app.POST("/project", handlers.SaveProjectHandler)

		app.GET("/clip", handlers.GetClipHandler)
		app.POST("/clip", handlers.SaveClipHandler)
	}
	router.Run(":" + os.Getenv(Port))
	// TODO RunTLS for logged-in stuff
}

func XHRMiddleware(ctx *gin.Context) {
	ctx.Set("is_xhr", ctx.Request.Header.Get("X-Requested-With") == "XMLHttpRequest")
}

func SessionMiddleware(r *models.Resource) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s, _ := models.InitSession(r, ctx)

		// TODO handle err

		ctx.Set("session", s)
	}
}

func staticPage(page string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, page+".tmpl.html", nil)
	}
}

func InGroup(minGroup int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session, _ := ctx.MustGet("session").(models.Session)

		// only All is allowed to be unauthenticated
		if !session.Auth && minGroup < All {
			if ctx.MustGet("is_xhr").(bool) {
				InternalError(ctx, NOT_LOGGED_IN, http.StatusForbidden)
			} else {
				// TODO last page url param to redir to cur page after auth
				ctx.Redirect(http.StatusSeeOther, "/auth")
			}
			ctx.Abort()
			return
		}

		if minGroup < session.Group {
			InternalError(ctx, INSUFFICIENT_PRIVS, http.StatusForbidden)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

const (
	Admin int = iota // privileged
	User             // regular logged-in user
	All              // any user, logged-in or logged-out
)
