package models

// a User model represents a person that can be logged in to the site
type User struct {
	// in schema
	Id    int64
	Email string // used as login name
	Pass  string // salted, hashed password

	// not in schema
	Res *Resource // ref to initialized resources
}

func InitUser(res *Resource, id int64) (u *User, err error) {
	u = &User{Res: res}
	err = u.Load(id)
	return
}

func (u *User) Load(id int64) (err error) {
	row := u.Res.LoadRowById("users", id)
	err = row.Scan(&u.Id, &u.Email, &u.Pass)
	return
}

func (u *User) LoadByEmail(email string) (err error) {
	row := u.Res.LoadRow("users", "email", email)
	err = row.Scan(&u.Id, &u.Email, &u.Pass)
	return
}

func (u *User) Store() (err error) {
	err = u.Res.StoreRow(
		"users",
		[]string{"email", "pass"},
		&u.Id,
		u.Email, u.Pass)

	return
}
