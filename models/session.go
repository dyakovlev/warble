package models

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/dyakovlev/warble/utils"
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
	s := Session{Res: r}

	if rawSessionId, badCookie := utils.GetSessionCookie(c); badCookie != nil && rawSessionId != "" {
		if decSessionId := r.Decid(rawSessionId); decSessionId != 0 {
			if noSessionRow := s.Load(decSessionId); noSessionRow != nil {
				s.UpdateSeen()
				return &s, nil
			}
		}
	}

	// TODO handle errors

	// if we're here we're making a new session

	s.Seen = time.Now()
	s.Group = 2 // All

	return &s, nil
}

func (s *Session) Load(id int) (err error) {
	row := s.Res.LoadRow("session", id)
	err = row.Scan(&s.Id, &s.Auth, &s.Group, &s.Seen, &s.Uid, &s.Pid)

	if err != nil {
		// TODO handle norows error
	}

	return err
}

func (s *Session) Store() error {
	_, err := s.Res.StoreRow(
		"session",
		[]string{"id", "auth", "grp", "seen", "uid", "pid"},
		&s.Id, &s.Auth, &s.Group, &s.Seen, &s.Uid, &s.Pid,
	)
	return err
}

func (s *Session) Expire() (err error) {
	s.Auth = false
	err = s.Store()

	if err != nil {
		// TODO log session expiration error
	}

	return err
}

func (s *Session) Authorize(uid int) error {
	s.Auth = true
	s.Uid = uid
	err := s.Store()

	if err != nil {
		s.Auth = false
		s.Uid = -1
	}

	return err
}

func (s *Session) Detach() error {
	s.Auth = false
	s.Uid = -1
	return s.Store()
}

func (s *Session) UpdateSeen() error {
	s.Seen = time.Now()
	return s.Store()
}
