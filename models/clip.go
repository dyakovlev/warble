package models

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// TODO how closely should this be integrated with the S3 clip storage?

type Clip struct {
	// in schema
	Id     int64
	S3Addr string // where to find the raw clip in s3

	// not in schema
	Res *Resource // ref to initialized resources
}

func InitClip(r *Resource, c *gin.Context) (*Clip, error) {
	c := Clip{Res: r}
	return c, nil
}

func (c *Clip) Load(id int64) (err error) {
	row := c.Res.LoadRowById("clips", id)
	err = row.Scan(&c.S3Addr)
	return
}

func (c *Clip) Store() (err error) {
	if c.Id == 0 {
		err = c.Res.DB.QueryRow("INSERT INTO clip (s3addr) VALUES ($1) RETURNING id",
			c.S3Addr).Scan(&c.Id)
	} else {
		_, err = c.Res.DB.Exec("UPDATE project SET s3addr=$1 WHERE id = $2",
			c.S3Addr, c.Id)
	}

	handleDBError("Project.Store", err)
}

func (c *Clip) String() string {
	return fmt.Sprintf("Clip{id:%v, S3Addr:%v, Resource:%v}",
		c.Id, c.S3Addr, c.Res)
}
