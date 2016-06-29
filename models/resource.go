package models

import (
	"bytes"
	"database/sql"
	"fmt"
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
	crypter := utils.NewIDCodec(crypterKey) // TODO globally initialize crypter, doesn't need to be a part of Resource

	return &Resource{db, crypter}, err
}

func (r *Resource) Encid(plain int64) string {
	return r.crypter.Encid(plain)
}

func (r *Resource) Decid(enc string) int64 {
	return r.crypter.Decid(enc)
}

func (r *Resource) LoadRowById(table string, id int64) *sql.Row {
	// TODO sanitize `table` param

	return r.db.QueryRow("SELECT * FROM $1 WHERE id=$2", table, id)
}

func (r *Resource) LoadRow(table string, col string, val string) *sql.Row {
	// TODO sanitize params

	return r.db.QueryRow("SELECT * FROM $1 WHERE $2=$3", table, col, val)
}

func (r *Resource) StoreRow(table string, fields []string, pkey *int64, params ...interface{}) (err error) {
	// TODO sanitize params

	if *pkey == 0 {
		err = r.db.QueryRow("INSERT INTO $1 VALUES ("+buildNumString(len(params))+") RETURNING id", buildParams(table, params)...).Scan(pkey)
	} else {
		_, err = r.db.Exec("UPDATE $1 SET ("+buildFieldString(fields)+") WHERE id=?", buildParams(table, params, pkey)...)
	}

	switch {
	case err != nil:
		utils.Error("[resource] StoreRow error:", err)
	}

	return
}

// returns "$2,$3,$4" if q=3
func buildNumString(q int) string {
	var out bytes.Buffer

	for i := 2; i < q+2; i++ {
		out.WriteString(fmt.Sprintf("$%d,", i))
	}

	outs := out.String()
	return strings.TrimRight(outs, ",")
}

// returns "field0=$2,field1=$3,.." for every field in fields
func buildFieldString(fields []string) string {
	var out bytes.Buffer

	for i := 2; i < len(fields)+2; i++ {
		out.WriteString(fmt.Sprintf("%s=$%d,", fields[i-2], i))
	}

	outs := out.String()
	return strings.TrimRight(outs, ",")
}

// combines passed params into a flat list
func buildParams(t string, p []interface{}, extra ...interface{}) []interface{} {
	return append(append([]interface{}{t}, p...), extra...)
}
