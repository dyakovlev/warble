package models

import (
	"time"

	"github.com/gin-gonic/gin"
)

// a Session model represents a logged-in session
type Session struct {
	// in schema
	id   int       // session id (primary key)
	auth bool      // has the user logged in
	grp  int       // what group does the logged-in user belong to
	seen time.Time // last seen
	uid  int       // associated user (if authenticated)
	pid  int       // associated project (last one worked on, for convenience)

	// not in schema
	resource *Resource // ref to initialized resources
}

func InitSession(r *Resource, c *gin.Context) (*Session, error) {
	if encSession, badCookie := c.Cookie(sessionCookie); badCookie != nil {
		s := Session{resource: r}
		if noSession := s.load(r.crypter.decid(sid)); noSession != nil {
			s.updateSeen()
			return &s, nil
		}
	}

	if badCookie || noSession {
		// TODO log
	}

	newSession := Session{group: All, seen: time.Now()}
	err := newSession.store()

	// set session cookie
	SetSessionCookie(c, r.crypter.Encid(newSession.id))

	return &newSession, nil
}

func (s *Session) load(id int) (err error) {
	if row, err := s.resource.LoadRow("session", id); err != nil {
		err = row.Scan(&s.id, &s.auth, &s.grp, &s.seen, &s.uid, &s.pid)
	}
}

func (s *Session) store() error {
	s.resource.StoreRow(
		"session",
		[]string{"id", "auth", "grp", "seen", "uid", "pid"},
		&s.id, &s.auth, &s.grp, &s.seen, &s.uid, &s.pid,
	)
}

func (s *Session) expire() (err error) {
	s.auth = false
	_, err = s.store()

	if err != nil {
		// TODO log session expiration error
	}
}

func (s *Session) authorize(uid int) (err error) {
	s.auth = true
	s.uid = id
	_, err = s.store()

	if err != nil {
		s.auth = false
		s.uid = nil
	}
}

func (s *Session) detach() (err error) {
	s.auth = false
	s.uid = nil
	_, err = s.store()
}

func (s *Session) updateSeen() (err error) {
	s.seen = time.Now()
	_, err = s.store()
}
