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
	if u.Id == 0 {
		err = u.Res.DB.QueryRow("INSERT INTO users (email, pass) VALUES ($1::varchar(255), $2::varchar(255)) RETURNING id",
			u.Email, u.Pass).Scan(&u.Id)
	} else {
		_, err = u.Res.DB.Exec("UPDATE users SET (email=$1,pass=$2) WHERE id=$3",
			u.Email, u.Pass, u.Id)
	}

	handleDBError(err)

	return
}
