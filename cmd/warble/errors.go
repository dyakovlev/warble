package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InternalError(ctx *gin.Context, err string, errCode int) {
	// TODO log the error

	ctx.String(errCode, http.StatusText(errCode)+": "+err)
}

const (
	INSUFFICIENT_PRIVS = "You can't do that."
	NOT_LOGGED_IN      = "You're not logged in."
)

var NoSuchUser error = errors.New("that's not a user that exists")
