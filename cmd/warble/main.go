package main

import (
	"os"
	"fmt"
	"database/sql"

	"github.com/gin-tonic/gin"
)

func main() {
	router := gin.Default()
	router.Static("/static", "static") // TODO eventually replace with CDN
	router.LoadHTMLGlob("templates/*.tmpl.html")

	db := sql.Open("postgres", os.Getenv("DB_URL"))
	ss := SessionStore{db}

	router.GET("/about",	 StaticPage("about"))
	router.GET("/status",	 WithSession(ss), WithGroup(Admin), StatusHandler)
	router.GET("/",			 WithSession(ss), RootHandler)

	app := router.Group("/", WithSession(ss), WithGroup(User))
	{
		router.GET("/config",	 GetConfigHandler)
		router.POST("/config",	 SaveConfigHandler)
		router.GET("/clip",		 GetClipHandler)
		router.POST("/clip",	 SaveClipHandler)
	}

	router.Run(":" + os.Getenv("PORT"))
	// RunTLS for logged-in stuff
}


func StaticPage(page string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, page + ".tmpl.html")
	}
}

// retrieves or generates a session
func WithSession(ss *SessionStore) gin.HandlerFunc {
	return func(ctx *Context) {
		// grab session token from request
		// see if it's in ss
		// load session into ctx if it is
		// make a new one if not

		c.Next()
	}
}


type Group int
const (
	Admin Group iota
	User
)

// generates a middleware that verifies request group fits passed param
func WithGroup(g Group) gin.HandlerFunc {
	return func(ctx *Context) {
		// get user permissions
		// if it satisfies g, pass
		// if not, abort with 500

		c.Next()
	}
}
