package main

import (
	"database/sql"
	_ "github.com/lib/pq"

	"github.com/gin-tonic/gin"
)

type Database struct {
	db *sql.DB
}

// middlewre that initializes supplied model and inserts into context
func (db *Database) WithModel(name string, initializer Initializer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		model, err := initializer(db, ctx.Request)
		// todo handle err

		ctx.Set(name, model)
		ctx.Next()
	}
}

func (db *Database) WithSession() gin.HandlerFunc {
	return db.withModel("session", InitSession)
}

func (db *Database) WithUser() gin.HandlerFunc {
	return db.withModel("user", InitUser)
}

func (db *Database) WithProject() gin.HandlerFunc {
	return db.withModel("project", InitProject)
}

func (db *Database) WithClip() gin.HandlerFunc {
	return db.withModel("clip", InitClip)
}

// shorthand type
type Initializer func(db *Database, r *http.Request)
