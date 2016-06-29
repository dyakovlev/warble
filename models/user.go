package models

// a User model represents a person that can be logged in to the site
type User struct {
	// in schema
	Id    int64
	Email string // used as login name
	Pass  string // salted, hashed password
	Admin bool   // is user an admin

	// not in schema
	Res *Resource // ref to initialized resources
}

func InitUser(s *Session) (u *User, err error) {
	u = &User{Res: s.Res}
	err = u.Load(s.Uid)
	return
}

func (u *User) Load(id int64) (err error) {
	row := u.Res.LoadRowById("users", id)
	err = row.Scan(&u.Email, &u.Pass, &u.Id, &u.Admin)
	return
}

func (u *User) LoadByEmail(email string) (err error) {
	row := u.Res.LoadRow("users", "email", email)
	err = row.Scan(&u.Email, &u.Pass, &u.Id, &u.Admin)
	return
}

func (u *User) Store() (err error) {
	if u.Id == 0 {
		err = u.Res.DB.QueryRow("INSERT INTO users (email, pass, admin) VALUES ($1::varchar(255), $2::varchar(255), $3::boolean) RETURNING id",
			u.Email, u.Pass, u.Admin).Scan(&u.Id)
	} else {
		_, err = u.Res.DB.Exec("UPDATE users SET email=$1, pass=$2, admin=$3 WHERE id=$3",
			u.Email, u.Pass, u.Admin, u.Id)
	}

	handleDBError("User.Store", err)

	return
}
