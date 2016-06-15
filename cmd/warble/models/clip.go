package main

import (
	"database/sql"
	"net/http"
)

// a Clip model represents some clip metadata (raw clips are stored in S3)
type Clip struct {
}

func (m Session) name() string { return "clip" }

func (m *Clip) Initialize(db *sql.DB, r *http.Request) *Clip {

}

func (m *Clip) GetById(id int) *Clip {

}

func (m *Clip) Upsert(r Row) (success bool) {

}
