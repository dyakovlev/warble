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

	auth := router.Group("/auth", sessionMiddleware(r))
	{
		auth.GET("", handlers.ServeAuthPage)
		auth.POST("/new", handlers.DoNewAccount)
		auth.POST("/login", handlers.DoAuth)
		auth.POST("/logout", handlers.DoLogout)
	}

	app := router.Group("/", sessionMiddleware(r), loggedIn)
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

func sessionMiddleware(r *models.Resource) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s, err := models.InitSession(r, ctx)

		if err != nil {
			utils.Error("[sessionMiddleware] couldn't initialize session")
			ctx.Abort()
		}

		ctx.Set("session", s)
	}
}

func staticPage(page string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, page+".tmpl.html", nil)
	}
}

func loggedIn(ctx *gin.Context) {
	session, _ := ctx.MustGet("session").(*models.Session)

	if !session.Auth {
		if ctx.MustGet("is_xhr").(bool) {
			ctx.JSON(http.StatusOK, gin.H{"error": "unauthorized"})
		} else {
			ctx.Redirect(http.StatusSeeOther, "/auth")
		}

		utils.Error("[loggedIn] aborting response due to not being authenticated")
		ctx.Abort()
	}
}
