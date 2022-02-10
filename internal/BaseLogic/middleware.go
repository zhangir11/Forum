package baselogic

import (
	"forum/pkg/funcs"
	"forum/pkg/sserver"
	"net/http"
)

//Middleware ...
func Middleware(f func(w http.ResponseWriter, r *http.Request, u *funcs.User), AuthCode int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		u, err := sserver.CheckIfUserAuthorized(r)

		if (AuthCode == 1 && (err != nil)) || (AuthCode == 2 && err == nil) {
			http.Redirect(w, r, "/main", http.StatusSeeOther)
			return
		}

		f(w, r, u)
	}
}
