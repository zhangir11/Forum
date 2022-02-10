package baselogic

import (
	"forum/pkg/funcs"
	"forum/pkg/sserver"
	"net/http"
)

//Logout ...
func Logout(w http.ResponseWriter, req *http.Request, u *funcs.User) {
	c, err := req.Cookie("authenticated")
	if err == nil {

		stmt, err := sserver.Db.Db.Prepare("DELETE FROM Sessions WHERE cookieValue = ?")
		if err != nil {
			return
		}

		stmt.Exec(c.Value)
		stmt.Close()
	}
	http.Redirect(w, req, "/main", http.StatusSeeOther)
}
