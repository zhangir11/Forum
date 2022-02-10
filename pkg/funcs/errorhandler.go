package funcs

import (
	"net/http"
	"text/template"
)

//ErrorHandler ...
func ErrorHandler(w http.ResponseWriter, status string, errorcase int) {
	tpl, err := template.ParseFiles("web/templates/error.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "500: Internal server error", http.StatusInternalServerError)
	} else {
		w.WriteHeader(errorcase)
		tpl.Execute(w, status)
	}
}
