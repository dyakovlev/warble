package main

import (
	"github.com/gin-tonic/gin"
	"net/http"
)

func InternalError(ctx *gin.Context, err string, errCode int) {
	// TODO log the error

	ctx.String(errCode, http.StatusText(errCode)+": "+err)
}

const (
	INSUFFICIENT_PRIVS = "You can't do that."
	NOT_LOGGED_IN      = "You're not logged in."
)
