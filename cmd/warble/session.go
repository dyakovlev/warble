package session

import (
	"github.com/gin-gonic/gin"
)

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
