package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProfileHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "user.tmpl.html", gin.H{})
}

func SaveProfileHandler(c *gin.Context) {

}
