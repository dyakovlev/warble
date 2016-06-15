package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

// SessionStore holds a bunch of sessions

type SessionStore struct {
	db *sql.DB			// keep a DB ref around
}

func (ss *SessionStore) Lookup(sid string) *Session {

}

func (ss *SessionStore) New() *Session {

}


type Session struct {
	sid string // session id
	groups []groups.Group // group session belongs to
	ctx *gin.Context // pointer to the request context
}

func Initialize(ctx *gin.Context) Session {
	// generate a session out of a context

}

func (this Session) String() string {
	fmt.Sprintf("session id: %s", this.sid)
	fmt.Sprintf("group: %s", this.group)
}
