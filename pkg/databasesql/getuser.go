package databasesql

import (
	"forum/pkg/funcs"
	"log"
)

// GetUserByID ...
func (d *DatabaseSQL) GetUserByID(uid int) (*funcs.User, error) {
	u := funcs.NewUser()
	if err := d.Db.QueryRow("SELECT firstname,lastname,username FROM Users WHERE id = ?", uid).Scan(
		&u.FirstName,
		&u.LastName,
		&u.UserName,
	); err != nil {
		return nil, err
	}
	u.ID = uid
	return u, nil
}

//GetUserByName ...
func (d *DatabaseSQL) GetUserByName(u *funcs.User) error {
	if err := d.Db.QueryRow("SELECT id, username,password FROM Users where username = ?", u.UserName).Scan(
		&u.ID,
		&u.UserName,
		&u.EncPwd,
	); err != nil {
		return err
	}
	return nil
}

//FindByEmail ...
func (d *DatabaseSQL) FindByEmail(u *funcs.User) error {
	if err := d.Db.QueryRow("SELECT id,email,password FROM Users where email = ?", u.Email).Scan(
		&u.ID,
		&u.Email,
		&u.EncPwd,
	); err != nil {
		log.Println(err, "Email err")
		return err
	}
	return nil
}

// GetUserByCookie ...
func (d *DatabaseSQL) GetUserByCookie(cookieValue string) (*funcs.User, error) {
	var userID int
	if err := d.Db.QueryRow("SELECT userID from Sessions WHERE cookieValue = ?", cookieValue).Scan(&userID); err != nil {
		return nil, err
	}
	u, err := d.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return u, nil
}
