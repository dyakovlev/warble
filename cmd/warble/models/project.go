package main

import (
	"database/sql"
	"net/http"
)

// a Project model represents a project configuration
type Project struct {
}

func (m Project) name() string { return "project" }

func (m Project) Initialize(db *sql.DB, r *http.Request) *Project {

}

func (m *Project) GetById(id int) *Project {

}

func (m *Project) Upsert(r Row) (success bool) {
}
