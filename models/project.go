package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// a Project model represents a project configuration (where a project is a song)
type Project struct {
	// in schema
	Id       int64
	Uid      int64     // userid of project owner
	Name     string    // username/unique-to-user-project-name
	JSON     string    // json blob project description
	Modified time.Time // last modified timestamp
	Public   boolean   // is the project publically visible

	// not in schema
	resource *Resource // ref to initialized resources
}

func InitProject(r *Resource, viewer int64, owner string, name string) (*Project, error) {
	p := Project{Res: r}

	err := p.LoadByName(fmt.Sprintf("%q/%q", owner, name))

	if err != nil {
		return nil, err
	} else if p.Uid != viewer && !p.Public {
		return nil, errors.New("can't view private projects that aren't your own")
	}

	return &p, nil
}

func (p *Project) Load(id int64) (err error) {
	row := p.Res.LoadRowById("project", id)
	err = row.Scan(&p.Id, &p.Uid, &p.Name, &p.JSON, &p.Modified, &p.Public)
	return
}

func (p *Project) LoadByName(name string) (err error) {
	row := p.Res.LoadRow("project", "name", name)
	err = row.Scan(&p.Id, &p.Uid, &p.Name, &p.JSON, &p.Modified, &p.Public)
	return
}

func (p *Project) Store() (err error) {
	if p.Id == 0 {
		err = p.Res.DB.QueryRow("INSERT INTO project (uid, name, json, mod, pub) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			p.Uid, p.Name, p.JSON, p.Modified, p.Public).Scan(&p.Id)
	} else {
		_, err = p.Res.DB.Exec("UPDATE project SET uid=$1, name=$2, json=$3, mod=$4, pub=$5  WHERE id = $6",
			p.Uid, p.Name, p.JSON, p.Modified, p.Public, p.Id)
	}

	handleDBError("Project.Store", err)
}

func (p *Project) String() string {
	return fmt.Sprintf("Project{id:%v, uid:%v, name:%v, mod:%v, pub:%v, Resource:%v}",
		p.Id, p.Uid, p.Name, p.JSON, p.Modified, p.Public, p.Res)
}
