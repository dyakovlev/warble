package main

import (
	"github.com/gin-tonic/gin"
)

// get project page
func GetProjectHandler(c *gin.Context) {
	session := c.MustGet("session")
	project := c.MustGet("project")
}

// save project info
func SaveProjectHandler(c *gin.Context) {
	session := c.MustGet("session")
	project := c.MustGet("project")

}

// retrieve clip info
func GetClipHandler(c *gin.Context) {
	session := c.MustGet("session")
	clip := c.MustGet("clip")
}

// save clip info
func SaveClipHandler(c *gin.Context) {
	session := c.MustGet("session")
	clip := c.MustGet("clip")

}
