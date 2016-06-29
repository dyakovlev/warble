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

	if rawSessionId, badCookie := utils.GetSessionCookie(c); badCookie != nil || rawSessionId == "" {
		utils.Error("InitSession: bad or no cookie:", rawSessionId, badCookie)
	} else if decSessionId := r.Decid(rawSessionId); decSessionId == 0 {
		utils.Error("InitSession: Decid couldn't decode session cookie", rawSessionId)
	} else if noSessionRow := s.Load(decSessionId); noSessionRow != nil {
		utils.Error("InitSession: couldn't load session", decSessionId, noSessionRow)
	} else {
		utils.Info("InitSession: successfully loaded session", decSessionId, "from cookie", rawSessionId)
		s.UpdateSeen()
		return &s, nil
	}

	// expire any existing session cookie if we've gone this far and it died..

	utils.ExpireSessionCookie(c)

	// if we're here we're making a new session (but not saving it until it gets authenticated)

	s.Seen = time.Now()
	s.Group = 2 // All

	utils.Info("InitSession: made a new session")

	return &s, nil
}

func (s *Session) Load(id int64) (err error) {
	row := s.Res.LoadRowById("session", id)
	err = row.Scan(&s.Auth, &s.Group, &s.Seen, &s.Uid, &s.Pid, &s.Id)

	if err != nil {
		// TODO handle norows error
	}

	return err
}

func (s *Session) Store() (err error) {
	if s.Id == 0 {
		err = s.Res.DB.QueryRow("INSERT INTO session (auth, grp, seen, uid, pid) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			s.Auth, s.Group, s.Seen, s.Uid, s.Pid).Scan(&s.Id)
	} else {
		_, err = s.Res.DB.Exec("UPDATE session SET auth=$1, grp=$2, seen=$3, uid=$4, pid=$5 WHERE id = $6",
			s.Auth, s.Group, s.Seen, s.Uid, s.Pid, s.Id)
	}

	handleDBError("Session.Store", err)

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

	if u.Admin {
		s.Group = 0
	} else {
		s.Group = 1
	}

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
	return fmt.Sprintf("Session{id:%v, uid:%v, pid:%v, auth:%v, group:%v, seen:%v, Resource:%v}",
		s.Id, s.Uid, s.Pid, s.Auth, s.Group, s.Seen, s.Res)
}
