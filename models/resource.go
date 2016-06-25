package models

import (
	"database/sql"
	_ "github.com/lib/pq"
	"strings"

	"github.com/dyakovlev/warble/utils"
)

// container for initialized resource singletons
type Resource struct {
	db      *sql.DB        // db connection
	crypter *utils.IDCodec // initialized id crypter (for public-facing IDs)
}

func NewResource(dbAddress string, crypterKey string) (*Resource, error) {
	db, err := sql.Open("postgres", dbAddress)
	crypter := utils.NewIDCodec(crypterKey)

	return &Resource{db, crypter}, err
}

func (r *Resource) Encid(plain int) string {
	return r.crypter.Encid(plain)
}

func (r *Resource) Decid(enc string) int {
	return r.crypter.Decid(enc)
}

func (r *Resource) LoadRow(table string, id int) *sql.Row {
	// TODO sanitize `table` param

	return r.db.QueryRow("SELECT * FROM ? WHERE id=?", table, id)
}

func (r *Resource) StoreRow(table string, fields []string, params ...interface{}) (int, error) {
	// TODO sanitize params

	var err error
	var res sql.Result
	pkey := params[0]

	// if supplied model doesn't have a defined primary key value, we make a new row and return
	// the id (which the model should then store)

	if pkey == nil {
		qs := strings.TrimRight(strings.Repeat("?,", len(params)), ",")
		res, err = r.db.Exec("INSERT INTO ? VALUES ("+qs+")", table, params...)
		pkey, err = res.LastInsertId()
		// TODO make sure pkey refs back into the model it's set from
	} else {
		fieldString := strings.Join(fields, "=?, ") + "=?"
		params = append(params, pkey) // fill out the id param
		_, err = r.db.Exec("UPDATE ? SET ("+fieldString+") WHERE id=?", table, params...)
	}

	// TODO what errors do we need to handle here

	switch {
	case err != nil:
		// log a fatal error
	}

	return pkey.(int), err
}
