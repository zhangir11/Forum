package funcs

import "net/http"

//User ...
type User struct {
	UserName  string
	FirstName string
	LastName  string
	ID        int
	Password  string
	Email     string
	EncPwd    string
}

//NewUser ...
func NewUser() *User {
	return &User{}
}

//Database ...
type Database interface {
	Initiate(string) error
	AddSession(http.ResponseWriter, string, *User) error
	GetUserByID(int) (*User, error)
	GetUserByName(*User) error
	PutUser(*User) error
	Close() error
	FindByEmail(*User) error
	GetUserByCookie(cookieValue string) (*User, error)
}
