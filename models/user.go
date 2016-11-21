package models

// a User model represents a person that can be logged in to the site
type User struct {
	// in schema
	Id    int64
	Name  string // unique name
	Email string // used for contact
	Pass  string // salted, hashed password

	// not in schema
	Res *Resource // ref to initialized resources
}

func InitUser(r *Resource, uid int64) (*User, error) {
	u = User{Res: r}
	err = u.Load(uid)
	return &u, err
}

func (u *User) Load(id int64) (err error) {
	row := u.Res.LoadRowById("users", id)
	err = row.Scan(&u.Name, &u.Email, &u.Pass, &u.Id)
	return
}

func (u *User) LoadByName(name string) (err error) {
	row := u.Res.LoadRow("users", "name", name)
	err = row.Scan(&u.Name, &u.Email, &u.Pass, &u.Id)
	return
}

func (u *User) LoadByEmail(email string) (err error) {
	row := u.Res.LoadRow("users", "email", email)
	err = row.Scan(&u.Name, &u.Email, &u.Pass, &u.Id)
	return
}

func (u *User) Store() (err error) {
	if u.Id == 0 {
		err = u.Res.DB.QueryRow("INSERT INTO users (name, email, pass) VALUES ($1, $2, $3) RETURNING id",
			u.Name, u.Email, u.Pass).Scan(&u.Id)
	} else {
		_, err = u.Res.DB.Exec("UPDATE users SET name=$1, email=$2, pass=$3 WHERE id=$4",
			u.Name, u.Email, u.Pass, u.Id)
	}

	handleDBError("User.Store", err)

	return
}

func (u *User) String() string {
	return fmt.Sprintf("User{id:%v, name:%v, email:%v, hashed pass:%v Resource:%v}",
		u.Id, u.Name, u.Email, u.Pass, u.Res)
}
