package models

import (
	"database/sql"
	_ "github.com/lib/pq"

	"github.com/dyakovlev/warble/utils"
)

// container for initialized resource singletons
type Resource struct {
	DB      *sql.DB        // db connection
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
	// sanitize table
	return r.DB.QueryRow("SELECT * FROM "+table+" WHERE id=$1", id)
}

func (r *Resource) LoadRow(table string, col string, val string) *sql.Row {
	// sanitize table
	return r.DB.QueryRow("SELECT * FROM "+table+" WHERE $1=$2", col, val)
}

func handleDBError(prefix string, err error) {
	switch {
	case err != nil:
		utils.Error(prefix, "DB error:", err)
	}
}
