package main

import (
	"net/http"
	"time"
)

// a Session model represents a logged-in session
type Session struct {
	// in schema
	id    int       // session id (primary key)
	auth  bool      // has the user logged in
	group int       // what group does the logged-in user belong to
	seen  time.Time // last seen
	uid   int       // associated user (if authenticated)
	pid   int       // associated project (last one worked on, for convenience)

	// not in schema
	resource *Resource // ref to initialized resources
}

func InitSession(res *Resource, req *http.Request) (*Session, error) {
	if rawCookie, noCookie := req.Cookie(SessionKey); noCookie != nil {
		if id, badCookie := decode(rawCookie); badCookie != nil {
			s := Session{resource: res}
			if noSession := s.load(id); noSession != nil {
				s.updateSeen()
				return &s, nil
			}
		}
	}

	if badCookie || noSession {
		// TODO log
	}

	newSession := Session{group: All, seen: time.Now()}
	newId, err := newSession.store()
	newSession.id = newId

	// TODO set the sid cookie

	return &newSession, nil
}

func (s *Session) load(id int) (err error) {
	if row, err := s.resource.loadRow("session", id); err != nil {
		err = row.Scan(&s.id, &s.auth, &s.group, &s.seen, &s.uid, &s.pid)
	}
}

func (s *Session) store() (id int, err error) {
	id, err := s.resource.storeRow(
		"session",
		[]string{"id", "auth", "group", "seen", "uid", "pid"},
		&s.id, &s.auth, &s.group, &s.seen, &s.uid, &s.pid,
	)
}

func (s *Session) expire() (err error) {
	s.auth = false
	_, err := s.store()
}

func (s *Session) authorize(id int) (err error) {
	s.auth = true
	s.id = id
	_, err := s.store()
}

func (s *Session) updateSeen() (err error) {
	s.seen = time.Now()
	_, err := s.store()
}
