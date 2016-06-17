package main

import (
	"net/http"
	"time"
)

// a Session model represents a logged-in session
type Session struct {
	// in schema
	sid   int       // session id (primary key)
	auth  bool      // has the user logged in
	group int       // what group does the logged-in user belong to
	seen  time.Time // last seen
	uid   int       // associated user (if authenticated)
	pid   int       // associated project (last one worked on, for convenience)

	// not in schema
	db *Database // ref to our database connection
}

func InitSession(db *Database, r *http.Request) (s *Session, err error) {
	if rawCookie, noCookie := r.Cookie(SessionKey); noCookie != nil {
		if id, badCookie := decode(rawCookie); badCookie != nil {
			s := Session{db: db}
			if noSession := s.load(id); noSession != nil {
				s.updateSeen()
				return &s, nil
			}
		}
	}

	if badCookie || noSession {
		// TODO log
	}

	newSession := Session{nil, false, All, time.Now()}
	newId := newSession.Store()

	// TODO set the sid cookie

	return &newSession, nil
}

func (s *Session) load(id int) (err error) {
	if row, err := s.db.loadRow("session", id); err != nil {
		err = row.Scan(&s.sid, &s.auth, &s.group, &s.seen, &s.uid, &s.pid)
	}
}

func (s *Session) store() (id int, err error) {
	id, err := s.db.storeRow("session", &s.sid, &s.auth, &s.group, &s.seen, &s.uid, &s.pid)
}

func (s *Session) expire() (err error) {
	s.auth = false
	_, err := s.store()
}

func (s *Session) authorize() (err error) {
	s.auth = true
	_, err := s.store()
}

func (s *Session) updateSeen() (err error) {
	s.seen = time.Now()
	_, err := s.store()
}
