package main

import (
	"database/sql"
	"net/http"
	"time"
)

// a Session model represents a logged-in session
type Session struct {
	sid           int  // session id (primary key)
	authenticated bool // has the user logged in
	group         int  // what group does the logged-in user belong to
	seen          Time // last seen
	uid           int  // associated user (if authenticated)
	pid           int  // associated project (last one worked on, for convenience)
}

func (m Session) name() string { return "session" }

func InitSession(db *sql.DB, r *http.Request) (session *Session, err error) {
	if rawCookie, noCookie := r.Cookie(SessionKey); noCookie != nil {
		if sid, badCookie := decode(rawCookie); badCookie != nil {
			if session, noSession := LoadSession(sid); noSession != nil {
				session.seen = time.Now()
				session.Upsert()
				return &session, nil
			}
		}
	}

	// if we got here something is fucky or there's no session, make a new one

	if badCookie || noSession {
		// TODO log
	}

	newSession := Session{nil, false, All, time.Now()}
	sid := newSession.Store()

	// TODO set the sid cookie

	return &newSession, nil
}

func (m *Session) Load(id int) *Session {

}

func (m *Session) Store() (sid int) {

}

func (m *Session) Expire() (success bool) {

}

func decode(token string) int {

}

func encode(token int) string {

}
