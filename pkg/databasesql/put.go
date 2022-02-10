package databasesql

import "forum/pkg/funcs"

//PutUser ...
func (d *DatabaseSQL) PutUser(u *funcs.User) error {

	s, err := d.Db.Prepare("INSERT INTO Users (firstname, lastname, username, email,password) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	res, err := s.Exec(u.FirstName, u.LastName, u.UserName, u.Email, u.EncPwd)
	if err != nil {
		return err
	}

	temp, _ := res.LastInsertId()
	u.ID = int(temp)
	err = s.Close()
	return err
}
