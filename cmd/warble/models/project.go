package main

import (
	"net/http"
)

// a Project model represents a project configuration (where a project is a song)
type Project struct {
	// in schema
	id       int
	owner    int       // userid of project creator
	p        string    // json blob project description
	modified time.Time // last modified timestamp

	// not in schema
	db *Database
}

func InitProject(db *Database, r *http.Request) (project *Project, err error) {

}

func (p *Project) load(db *Database, id int) (err error) {

}

func (p *Project) store() (id int) {
}
