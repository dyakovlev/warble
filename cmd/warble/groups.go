package groups

import (
	"session"
)

const All // admit all requests
const Admin // admit requests with an admin session
const LoggedIn // admit requests with any active session

type Group

// depending on the group type, pick a stragety for validating the session into that group
func Admit(group Group, s Session) (ok bool){

}
