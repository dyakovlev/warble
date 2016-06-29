package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dyakovlev/warble/models"
)

func GetProfileHandler(c *gin.Context) {
	session, _ := c.MustGet("session").(*models.Session)
	user, _ := models.InitUser(session)

	c.HTML(http.StatusOK, "user.tmpl.html", gin.H{
		"email": user.Email})
}

func SaveProfileHandler(c *gin.Context) {

}
