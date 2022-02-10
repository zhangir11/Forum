package delrate

import (
	"forum/pkg/funcs"
	"net/http"
)

func DeleteSmth(w http.ResponseWriter, r *http.Request, u *funcs.User) {
	funcs.ErrorHandler(w, "Not Found", http.StatusNotFound)
}
