package mypostpage

import (
	"forum/internal/post"
	"forum/pkg/funcs"
	"log"
	"net/http"
	"text/template"
)

func MyPostPage(w http.ResponseWriter, r *http.Request, u *funcs.User) {
	templates := []string{
		"./web/templates/somemain.html",
		"./web/templates/somenavbar.html",
		"./web/templates/somefooter.html",
	}

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		log.Println(err)
		funcs.ErrorHandler(w, "InternalServer error", http.StatusInternalServerError)
		return
	}

	if u == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	posts := post.GetPostsByAuthor(u)
	data := Data{
		Posts:         posts,
		Authenticated: true,
		CurrentUrl:    r.URL.Path,
		Page:          "My Posts",
		Username:      u.UserName,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
		funcs.ErrorHandler(w, "InternalServer error", http.StatusInternalServerError)
	}
}
