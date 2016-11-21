package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/dyakovlev/warble/models"
)

const ProjectTemplate = "content-project.tmpl.html"

func ServeProjectPage(c *gin.Context) {
	session, _ := c.MustGet("session").(*models.Session)

	username := c.Param("name")
	projectname := c.Param("project")

	p, err := models.InitProject(session.Res, session.Uid, username, projectname)

	if err != nil {
		c.HTML(http.StatusOK, ErrorTemplate, gin.H{})
	}

	c.HTML(http.StatusOK, ProjectTemplate, gin.H{
		"project": p,
	})
}

func SaveProjectHandler(c *gin.Context) {

}

func GetClipHandler(c *gin.Context) {

}

func SaveClipHandler(c *gin.Context) {

}
