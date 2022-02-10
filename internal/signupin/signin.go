package signupin

import (
	"forum/pkg/funcs"
	"forum/pkg/sserver"
	"html/template"
	"log"
	"net/http"
	"strings"
)

//Login ...
func Login(w http.ResponseWriter, r *http.Request, u *funcs.User) {
	data := &ErrorsSignUp{Page: "login"}

	if r.URL.Path != "/login" {
		funcs.ErrorHandler(w, "404 Not Found", http.StatusNotFound)
		return
	}

	tpl, err := template.ParseFiles("web/templates/login.html")
	if err != nil {
		log.Println("Failed on Login", err.Error())
		funcs.ErrorHandler(w, "Something went wrong", 500)
		return
	}

	if r.Method == "GET" {
		if err = tpl.Execute(w, data); err != nil {
			log.Println(err.Error() + " : signin.go")
			funcs.ErrorHandler(w, "Internal Error", http.StatusInternalServerError)
		}
		return
	}

	u.UserName = r.PostFormValue("Username")
	u.Password = r.PostFormValue("Password")
	unameOrEmail := UnameOrEmail(u.UserName)

	if unameOrEmail {
		err = sserver.Db.FindByEmail(u)
	} else {
		err = sserver.Db.GetUserByName(u)
	}

	if err != nil {
		data.Email = "Incorrect login or password"
		if err := tpl.Execute(w, data); err != nil {
			log.Println("Failed on SignUp", err.Error())
		}
		return
	}

	auth := funcs.ComparePassword(u.EncPwd, u.Password)
	if !auth {
		data.Email = "Incorrect login or password"
		if err := tpl.Execute(w, data); err != nil {
			log.Println("Failed on SignUp", err.Error())
		}
		return
	}

	if err := sserver.Servak.Database.AddSession(w, "authenticated", u); err != nil {
		funcs.ErrorHandler(w, "Something went wrong", http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	http.Redirect(w, r, "/main", http.StatusSeeOther)
}

// UnameOrEmail -->
func UnameOrEmail(query string) bool {
	if ok := strings.Contains(query, "@"); !ok {
		return false
	}
	return true
}
