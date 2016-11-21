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
	Id   int64     // session id (primary key)
	Uid  int64     // associated user (if authenticated)
	Pid  int64     // associated project (last one worked on, for convenience)
	Auth bool      // has the user logged in
	Seen time.Time // last seen

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

	// make a new session but don't save it (we don't need a row in the DB for every anon visit)
	s.Seen = time.Now()

	utils.Info("InitSession: making a new session", s)

	return &s, nil
}

func (s *Session) Load(id int64) (err error) {
	row := s.Res.LoadRowById("session", id)
	err = row.Scan(&s.Id, &s.Uid, &s.Pid, &s.Auth, &s.Seen)

	if err != nil {
		// TODO handle norows error
	}

	return err
}

func (s *Session) Store() (err error) {
	if s.Id == 0 {
		err = s.Res.DB.QueryRow("INSERT INTO session (uid, pid, auth, seen) VALUES ($1, $2, $3, $4) RETURNING id",
			s.Uid, s.Pid, s.Auth, s.Seen).Scan(&s.Id)
	} else {
		_, err = s.Res.DB.Exec("UPDATE session SET uid=$1, pid=$2, auth=$3, seen=$4,  WHERE id = $5",
			s.Uid, s.Pid, s.Auth, s.Seen, s.Id)
	}

	handleDBError("Session.Store", err)

	return
}

func (s *Session) Expire() (err error) {
	s.Auth = false
	err = s.Store()

	if err != nil {
		// TODO handle session expiration error
	}

	return err
}

func (s *Session) Authorize(u *User) error {
	s.Auth = true
	s.Uid = u.Id

	err := s.Store()

	if err != nil {
		s.Auth = false
		s.Uid = 0
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
	return fmt.Sprintf("Session{id:%v, uid:%v, pid:%v, auth:%v, seen:%v, Resource:%v}",
		s.Id, s.Uid, s.Pid, s.Auth, s.Seen, s.Res)
}
