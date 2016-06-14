package groups

import (
	"session"
)

/* group and user permission management utils */

type Group struct {
	id int
	name string
}

const All = Group{0, "Everyone"}
const Admin = Group{1, "Admins"}
const LoggedIn = Group{2, "Logged-in Users"}

func (this Group) String() string {
	fmt.Sprintf(this.name)
}

// depending on the group type, pick a stragety for validating the session into that group
func Admit(group Group, s session.Session) (ok bool){

}
