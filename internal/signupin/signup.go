package signupin

import (
	"forum/pkg/funcs"
	"forum/pkg/sserver"
	"html/template"
	"log"
	"net/http"
)

//SignUp ...
func SignUp(res http.ResponseWriter, req *http.Request, u *funcs.User) {
	tpl, err := template.ParseFiles("web/templates/login.html")
	if err != nil {
		log.Println("Failed on SignUp", err.Error())
		funcs.ErrorHandler(res, "Something went wrong", 500)
		return
	}
	success := true
	errorsSignUp := &ErrorsSignUp{Page: "signup"}

	if req.Method == http.MethodGet {
		tpl.Execute(res, errorsSignUp)
		return
	}

	u.FirstName = req.PostFormValue("Firstname")
	u.LastName = req.PostFormValue("Lastname")
	u.UserName = req.PostFormValue("Username")
	u.Email = req.PostFormValue("Email")
	u.Password = req.PostFormValue("Password")
	ConfirmPwd := req.PostFormValue("Confirm")
	var matchL, matchF bool = true, true

	// Validity of Email
	matchEmail := rxEmail.Match([]byte(u.Email))
	if !matchEmail {
		errorsSignUp.Email = invalidEmail
		success = false
	}

	// Username
	matchUname := rxUname.Match([]byte(u.UserName))
	if !matchUname {
		errorsSignUp.Username = invalidUsername
		success = false
	}

	// Password
	matchPwd := rxPwd.Match([]byte(u.Password))
	if !matchPwd {
		success = false
		errorsSignUp.Password = invalidPass
	}

	if u.Password != ConfirmPwd {
		errorsSignUp.Confirm = invalidConfirm
		success = false
	}
	// First and Last name
	if u.FirstName != "" && u.LastName != "" {
		matchF = rxUname.Match([]byte(u.FirstName))
		matchL = rxUname.Match([]byte(u.LastName))

		if !matchF || !matchL {
			errorsSignUp.FLName = invalidFLName
			success = false
		}
	} else {
		success = false
	}

	if success {
		encryptPass := funcs.HashPassword(u.Password)
		u.EncPwd = encryptPass                    // fill with Encrypted Password
		err := sserver.Servak.Database.PutUser(u) // Sending

		if err != nil {
			log.Println(err.Error() + " Signig Up")
			success = false
			errorsSignUp.Email = otherEmail
		}
	}

	if success {
		if err := sserver.Servak.Database.AddSession(res, "authenticated", u); err != nil {
			funcs.ErrorHandler(res, "Something went wrong", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		http.Redirect(res, req, "/main", http.StatusSeeOther)
		return
	}

	if err := tpl.Execute(res, errorsSignUp); err != nil {
		log.Println("Failed on SignUp", err.Error())
	}
}
