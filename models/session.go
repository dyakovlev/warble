package models

import (
	"time"

	"github.com/gin-gonic/gin"
)

// a Session model represents a logged-in session
type Session struct {
	// in schema
	Id    int       // session id (primary key)
	Auth  bool      // has the user logged in
	Group int       // what group does the logged-in user belong to
	Seen  time.Time // last seen
	Uid   int       // associated user (if authenticated)
	Pid   int       // associated project (last one worked on, for convenience)

	// not in schema
	Res *Resource // ref to initialized resources
}

func InitSession(r *Resource, c *gin.Context) (*Session, error) {
	if encSession, badCookie := GetSessionCookie(c); badCookie != nil {
		s := Session{Res: r}
		if noSession := s.Load(r.Decid(sid)); noSession != nil {
			s.UpdateSeen()
			return &s, nil
		}
	}

	if badCookie || noSession {
		// TODO log
	}

	newSession := Session{Auth: false, Group: All, Seen: time.Now()}
	return &newSession, nil
}

func (s *Session) Load(id int) (err error) {
	if row, err := s.Res.LoadRow("session", id); err != nil {
		err = row.Scan(&s.Id, &s.Auth, &s.Grp, &s.Seen, &s.Uid, &s.Pid)
	}
}

func (s *Session) Store() error {
	s.Res.StoreRow(
		"session",
		[]string{"id", "auth", "group", "seen", "uid", "pid"},
		&s.Id, &s.Auth, &s.Group, &s.Seen, &s.Uid, &s.Pid,
	)
}

func (s *Session) Expire() (err error) {
	s.Auth = false
	_, err = s.Store()

	if err != nil {
		// TODO log session expiration error
	}
}

func (s *Session) Authorize(uid int) (err error) {
	s.Auth = true
	s.Uid = id
	_, err = s.Store()

	if err != nil {
		s.Auth = false
		s.Uid = nil
	}
}

func (s *Session) Detach() (err error) {
	s.Auth = false
	s.Uid = nil
	_, err = s.Store()
}

func (s *Session) UpdateSeen() (err error) {
	s.seen = time.Now()
	_, err = s.Store()
}
