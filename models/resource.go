package models

import (
	"database/sql"
	_ "github.com/lib/pq"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/dyakovlev/warble/utils"
)

type Initializer func(r *Resource, c *gin.Context)

// container for initialized resource singletons
type Resource struct {
	db      *sql.DB        // db connection
	crypter *utils.IDCodec // initialized id crypter (for public-facing IDs)
}

func NewResource(dbAddress string, crypterKey string) (*Resource, err) {
	r := Resource{}

	r.db, err = sql.Open("postgres", dbAddress)
	r.crypter = utils.NewIDCodec(crypterKey)

	return &r, err
}

func (r *Resource) withModel(name string, initializer Initializer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		model, err := initializer(r, ctx)

		// TODO handle err

		ctx.Set(name, model)
		ctx.Next()
	}
}

func (r *Resource) WithSession() gin.HandlerFunc {
	return r.withModel("session", InitSession)
}

func (r *Resource) WithUser() gin.HandlerFunc {
	return r.withModel("user", InitUser)
}

func (r *Resource) WithProject() gin.HandlerFunc {
	return r.withModel("project", InitProject)
}

func (r *Resource) WithClip() gin.HandlerFunc {
	return r.withModel("clip", InitClip)
}

func (r *Resource) Encid(plain int) string {
	return r.crypter.Encid(plain)
}

func (r *Resource) Decid(enc string) int {
	return r.crypter.Decid(enc)
}

func (r *Resource) LoadRow(table string, id int) (*sql.Row, error) {
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

func (r *Resource) StoreRow(table string, fields []string, params ...interface{}) (int, error) {
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
