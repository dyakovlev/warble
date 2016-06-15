package main

import (
	"net/http"
	"time"
)

// a Session model represents a logged-in session
type Session struct {
	sid   int  // session id (primary key)
	auth  bool // has the user logged in
	group int  // what group does the logged-in user belong to
	seen  Time // last seen
	uid   int  // associated user (if authenticated)
	pid   int  // associated project (last one worked on, for convenience)
}

func InitSession(db *Database, r *http.Request) (session *Session, err error) {
	if rawCookie, noCookie := r.Cookie(SessionKey); noCookie != nil {
		if id, badCookie := decode(rawCookie); badCookie != nil {
			if session, noSession := LoadSession(db, id); noSession != nil {
				session.seen = time.Now()
				session.Store(db)

				return &session, nil
			}
		}
	}

	// if we got here something is fucky or there's no session, make a new one

	if badCookie || noSession {
		// TODO log
	}

	newSession := Session{nil, false, All, time.Now()}
	newId := newSession.Store()

	// TODO set the sid cookie

	return &newSession, nil
}

func LoadSession(db *Database, id int) (session *Session, err error) {
	if r, err := db.Load("session", id); err != nil {
		return &Session{r.sid, r.auth, r.group, r.seen, r.uid, r.pid}, nil
	}
}

func (m *Session) Store() (id int) {

}

func (m *Session) Expire() (success bool) {

}

func decode(token string) int {

}

func encode(token int) string {

}
