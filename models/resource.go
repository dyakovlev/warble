package models

import (
	"database/sql"
	"fmt"
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
	if err != nil {
		utils.Error("[NewResource] failed to create DB connection", err)
		return nil, err
	}

	// Open doesn't actually touch the DB, need to ping to see if we can talk to it
	err = db.Ping()
	if err != nil {
		utils.Error("[NewResource] failed to ping DB", err)
		return nil, err
	}

	crypter, err := utils.NewIDCodec(crypterKey)
	if err != nil {
		utils.Error("[NewResource] error calling NewIDCodec", err)
		return nil, err
	}

	return &Resource{db, crypter}, nil
}

func (r *Resource) Encid(plain int64) string {
	return r.crypter.Encid(plain)
}

func (r *Resource) Decid(enc string) int64 {
	return r.crypter.Decid(enc)
}

func (r *Resource) LoadRowById(table string, id int64) *sql.Row {
	// TODO sanitize table
	return r.DB.QueryRow("SELECT * FROM "+table+" WHERE id=$1", id)
}

func (r *Resource) LoadRow(table string, col string, val string) *sql.Row {
	// TODO sanitize table
	return r.DB.QueryRow("SELECT * FROM "+table+" WHERE $1=$2", col, val)
}

func handleDBError(prefix string, err error) {
	switch {
	case err != nil:
		utils.Error(fmt.Sprintf("[%v] DB error: %v", prefix, err))
	}
}
