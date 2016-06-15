package main

import (
	"net/http"
)

// a Clip model represents some clip metadata (raw clips are stored in S3)
type Clip struct {
}

func InitClip(db *Database, r *http.Request) (clip *Clip, err error) {

}

func LoadClip(db *Database, id int) (clip *Clip, err error) {

}

func (m *Clip) Store() (id int) {

}
