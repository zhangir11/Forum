package signupin

import "regexp"

//ErrorsSignUp ...
type ErrorsSignUp struct {
	Page, FLName, Username, Email, Password, Confirm string
}

// RegExp patterns to validate user's input
var (
	invalidEmail    = "Please enter a valid e-mail address, e.g. yourmail@example.com"
	invalidUsername = "Please use latin alphabet, lowercase and uppercase characters with numbers or special characters"
	invalidPass     = "The password must include one lowercase, uppercase, special characters and number"
	invalidFLName   = "Please use latin alphabet [a-z, A-Z]"
	invalidConfirm  = "The password confirmation does not match"
	otherEmail      = "Email or Username has alraedy taken"
	rxEmail         = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	rxUname         = regexp.MustCompile(`^[\w]`)
	rxPwd           = regexp.MustCompile(`^[#\w@-]{8,20}`)
)
