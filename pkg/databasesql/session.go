package databasesql

import (
	"forum/pkg/funcs"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

// AddSession ...
func (s *DatabaseSQL) AddSession(w http.ResponseWriter, sessionName string, u *funcs.User) error {
	// var s controller.Server
	cookieSession := &http.Cookie{
		Name:     sessionName,
		Value:    uuid.NewV4().String(),
		MaxAge:   900,
		HttpOnly: true,
	}

	http.SetCookie(w, cookieSession)

	stmnt, err := s.Db.Prepare("DELETE from Sessions WHERE userID = ? ")
	if err != nil {
		return err
	}

	_, err = stmnt.Exec(u.ID)
	if err != nil {
		return err
	}

	stmnt, err = s.Db.Prepare("INSERT INTO Sessions (userID, cookieName, cookieValue) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmnt.Exec(u.ID, cookieSession.Name, cookieSession.Value)
	if err != nil {
		return err
	}

	stmnt.Close()
	return nil

}
