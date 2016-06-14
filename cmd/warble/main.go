package main

import (
	"routes"
	"session"

	"github.com/gin-tonic/gin"
)

func main() {
	// grab config from environment

	// make conns to storage

	// init server
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")

	// assets
	router.Static("/static", "static")

	// routes
	router.GET("/about", serveStatic("about.tmpl.html"))
	router.GET("/", authorize(groups.All, routes.GetApp))
	router.GET("/config", authorize(groups.LoggedIn, routes.GetConfig))
	router.POST("/config", authorize(groups.LoggedIn, routes.SaveConfig))
	router.GET("/clip", authorize(groups.LoggedIn, routes.GetClip))
	router.POST("/clip", authorize(groups.LoggedIn, routes.SaveClip))
	router.GET("/status", authorize(groups.Admin, routes.Status))

	// serve
	router.Run()
}


func serveStatic(template string) func(ctx *gin.Context){
	return func(ctx *gin.Context) { ctx.HTML(http.StatusOK, template) }
}

func authorize(group groups.Group, wrappedRoute func(ctx *gin.Context)) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		if ok, s := session.Initialize(ctx); ok {

		} else {
			return
		}

		if groups.Admit(group, s) {

		} else {
			return
		}

		if ok, out := wrappedRoute(s); ok {

		} else {
			return
		}
	}
}

