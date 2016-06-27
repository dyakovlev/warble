package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// get user profile page
func GetUserHandler(c *gin.Context) {

	// get list of all projects

	c.HTML(http.StatusOK, "user.tmpl.html", gin.H{})
}

// save profile info
func SaveUserHandler(c *gin.Context) {

}
