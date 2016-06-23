package main

import (
	"net/http"
	"time"
)

// a Project model represents a project configuration (where a project is a song)
type Project struct {
	// in schema
	id       int
	owner    int       // userid of project creator
	p        string    // json blob project description
	modified time.Time // last modified timestamp

	// not in schema
	resource *Resource // ref to initialized resources
}

func InitProject(res *Resource, req *http.Request) (project *Project, err error) {

}

func (p *Project) load(id int) (err error) {

}

func (p *Project) store() (id int) {
}
