package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dyakovlev/warble/models"
)

const OwnProfileTemplate = "content-profile.tmpl.html"
const OtherProfileTemplate = "content-profile-external.tmpl.html"

func GetProfileHandler(c *gin.Context) {
	session, _ := c.MustGet("session").(*models.Session)

	user = models.User{Res: r}
	err := user.LoadByName(c.Param("name"))

	if err != nil {
		c.HTML(http.StatusNotFound, ErrorTemplate, gin.H{})
	}

	if user.Id == session.Uid {
		c.HTML(http.StatusOK, OwnProfileTemplate, gin.H{
			"email": user.Email})
	} else {
		c.HTML(http.StatusOK, OtherProfileTemplate, gin.H{})
	}
}

func SaveProfileHandler(c *gin.Context) {

}
