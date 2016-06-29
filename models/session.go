package models

import (
	"fmt"
	"time"

	"github.com/dyakovlev/warble/utils"
	"github.com/gin-gonic/gin"
)

// a Session model represents a logged-in session
type Session struct {
	// in schema
	Id    int64     // session id (primary key)
	Auth  bool      // has the user logged in
	Group int       // what group does the logged-in user belong to
	Seen  time.Time // last seen
	Uid   int64     // associated user (if authenticated)
	Pid   int64     // associated project (last one worked on, for convenience)

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

	// if we're here we're making a new session (but not saving it until it gets authenticated)

	s.Seen = time.Now()
	s.Group = 2 // All

	return &s, nil
}

func (s *Session) Load(id int64) (err error) {
	row := s.Res.LoadRowById("session", id)
	err = row.Scan(&s.Id, &s.Auth, &s.Group, &s.Seen, &s.Uid, &s.Pid)

	if err != nil {
		// TODO handle norows error
	}

	return err
}

func (s *Session) Store() (err error) {
	if s.Id == 0 {
		err = s.Res.DB.QueryRow("INSERT INTO session (auth, grp, seen, uid, pid) VALUES ($1::boolean, $2::integer, $3::timestamp without time zone, $4::integer, $5::integer) RETURNING id",
			s.Auth, s.Group, s.Seen, s.Uid, s.Pid).Scan(&s.Id)
	} else {
		_, err = s.Res.DB.Exec("UPDATE session SET (auth=$1,grp=$2,seen=$3,uid=$4,pid=$5) WHERE id=$6",
			s.Auth, s.Group, s.Seen, s.Uid, s.Pid, s.Id)
	}

	handleDBError(err)

	return
}

func (s *Session) Expire() (err error) {
	s.Auth = false
	err = s.Store()

	if err != nil {
		// TODO log session expiration error
	}

	return err
}

func (s *Session) Authorize(u *User) error {
	s.Auth = true
	s.Uid = u.Id
	s.Group = u.Group

	err := s.Store()

	if err != nil {
		s.Auth = false
		s.Uid = 0
		s.Group = 2
	}

	return err
}

func (s *Session) Detach() error {
	s.Auth = false
	s.Uid = 0
	return s.Store()
}

func (s *Session) UpdateSeen() error {
	s.Seen = time.Now()
	return s.Store()
}

func (s *Session) String() string {
	return fmt.Sprintf("session id:%v\nuid:%v\npid:%v\nauth:%v\nseen:%v\nResource:%v",
		s.Id, s.Uid, s.Pid, s.Auth, s.Seen, s.Res)
}
