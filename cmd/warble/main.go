package main

import (
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

	if err != nil {
		panic(err)
	}

	router.GET("/", staticPage("index"))

	router.GET("/about", staticPage("about"))

	// TODO basic auth
	router.GET("/status", handlers.StatusHandler)

	auth := router.Group("/auth", SessionMiddleware(r))
	{
		auth.GET("", handlers.GetAuthHandler)
		auth.POST("/new", handlers.DoNewAccountHandler)
		auth.POST("/login", handlers.DoAuthHandler)
		auth.POST("/logout", handlers.DoLogoutHandler)
	}

	app := router.Group("/", SessionMiddleware(r), LoggedIn)
	{
		app.GET("/profile", handlers.GetProfileHandler)
		app.POST("/profile", handlers.SaveProfileHandler)

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

func LoggedIn(ctx *gin.Context) {
	session, _ := ctx.MustGet("session").(*models.Session)
	if !session.Auth {
		if ctx.MustGet("is_xhr").(bool) {
			ctx.JSON(http.StatusOK, gin.H{"error": "unauthorized"})
		} else {
			// TODO last page url param to redir to cur page after auth
			ctx.Redirect(http.StatusSeeOther, "/auth")
		}

		utils.Error("[LoggedIn] aborting response due to not being authenticated")
		ctx.Abort()
	}
}
