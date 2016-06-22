package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"strings"

	"github.com/gin-tonic/gin"
)

type Initializer func(db *Database, r *http.Request)

type Database struct {
	DB      *sql.DB    // db connection
	crypter *IDCrypter // initialized id crypter (for public-facing IDs)
}

func (db *Database) withModel(name string, initializer Initializer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		model, err := initializer(db, ctx.Request)
		// todo handle err

		ctx.Set(name, model)
		ctx.Next()
	}
}

func (db *Database) withSession() gin.HandlerFunc {
	return db.withModel("session", InitSession)
}

func (db *Database) withUser() gin.HandlerFunc {
	return db.withModel("user", InitUser)
}

func (db *Database) withProject() gin.HandlerFunc {
	return db.withModel("project", InitProject)
}

func (db *Database) withClip() gin.HandlerFunc {
	return db.withModel("clip", InitClip)
}

func (db *Database) loadRow(table string, id int) (res *sql.Row, err error) {
	// TODO sanitize `table` param

	res, err := db.DB.QueryRow("SELECT * FROM ? WHERE id=?", table, id)

	switch {
	case err == sql.ErrNoRows:
		// log an attempt to load a row that doesn't exist
	case err != nil:
		// log a fatal error
	}

	return res, err
}

func (db *Database) storeRow(table string, fields ...string, params ...interface{}) (pkey int, err error) {
	// TODO sanitize params

	pkey := params[0]

	if pkey != nil {
		qs := strings.TrimRight(strings.Repeat("?,", len(params)), ",")
		pkey, err := db.DB.Query("INSERT INTO ? VALUES ("+qs+")", table, params...)
	} else {
		fieldString := strings.Join(fields, "=?, ") + "=?"
		append(params, pkey) // fill out the id param
		_, err := db.DB.Query("UPDATE ? SET ("+fieldString+") WHERE id=?", table, params...)
	}

	switch {
	case err != nil:
		// log a fatal error
	}

	return pkey, err
}
