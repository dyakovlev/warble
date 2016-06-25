package models

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// a Clip model represents some clip metadata (raw clips are stored in S3)
type Clip struct {
	// in schema
	id       int
	s3_addr  string    // where to find the raw clip in s3
	modified time.Time // last modified timestamp

	// not in schema
	resource *Resource // ref to initialized resources
}

func InitClip(r *Resource, c *gin.Context) (c *Clip, err error) {

}

func (c *Clip) load(id int) (err error) {

}

func (m *Clip) store() (id int, err error) {

}
