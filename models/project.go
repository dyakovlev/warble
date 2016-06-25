package models

import (
	"time"

	"github.com/gin-gonic/gin"
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

func InitProject(r *Resource, c *gin.Context) (project *Project, err error) {
	return &Project{}, nil
}

func (p *Project) load(id int) (err error) {
	return nil
}

func (p *Project) store() (id int, err error) {
	return 0, nil
}
