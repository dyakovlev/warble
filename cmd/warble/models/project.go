package main

import (
	"net/http"
)

// a Project model represents a project configuration
type Project struct {
}

func InitProject(db *Database, r *http.Request) (project *Project, err error) {

}

func LoadProject(db *Database, id int) (project *Project, err error) {

}

func (m *Project) Store() (id int) {
}
