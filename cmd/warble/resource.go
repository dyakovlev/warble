package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"strings"

	"github.com/gin-tonic/gin"
)

type Initializer func(db *Database, r *http.Request)

// container for initialized resource singletons
type Resource struct {
	db      *sql.DB    // db connection
	crypter *IDCrypter // initialized id crypter (for public-facing IDs)
}

func (r *Resource) withModel(name string, initializer Initializer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		model, err := initializer(r, ctx.Request)
		// todo handle err

		ctx.Set(name, model)
		ctx.Next()
	}
}

func (r *Resource) withSession() gin.HandlerFunc {
	return r.withModel("session", InitSession)
}

func (r *Resource) withUser() gin.HandlerFunc {
	return r.withModel("user", InitUser)
}

func (r *Resource) withProject() gin.HandlerFunc {
	return r.withModel("project", InitProject)
}

func (r *Resource) withClip() gin.HandlerFunc {
	return r.withModel("clip", InitClip)
}

func (r *Resource) loadRow(table string, id int) (res *sql.Row, err error) {
	// TODO sanitize `table` param

	res, err := r.db.QueryRow("SELECT * FROM ? WHERE id=?", table, id)

	switch {
	case err == sql.ErrNoRows:
		// log an attempt to load a row that doesn't exist
	case err != nil:
		// log a fatal error
	}

	return res, err
}

func (r *Resource) storeRow(table string, fields []string, params ...interface{}) (pkey int, err error) {
	// TODO sanitize params

	pkey := params[0]

	// if supplied model doesn't have a defined primary key value, we make a new row and return
	// the id (which the model should then store)

	if pkey == nil {
		qs := strings.TrimRight(strings.Repeat("?,", len(params)), ",")
		res, err := r.db.Exec("INSERT INTO ? VALUES ("+qs+")", table, params...)
		pkey = res.LastInsertId()
		// TODO make sure pkey refs back into the model it's set from
	} else {
		fieldString := strings.Join(fields, "=?, ") + "=?"
		append(params, pkey) // fill out the id param
		_, err := r.db.Exec("UPDATE ? SET ("+fieldString+") WHERE id=?", table, params...)
	}

	// TODO what errors do we need to handle here

	switch {
	case err != nil:
		// log a fatal error
	}

	return pkey, err
}
