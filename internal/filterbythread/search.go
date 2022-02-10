package filterbythread

import (
	"forum/internal/likedposts"
	"forum/internal/mypostpage"
	"forum/internal/post"
	"forum/pkg/funcs"
	"log"
	"net/http"
	"text/template"
)

func FilterByThread(w http.ResponseWriter, r *http.Request, u *funcs.User) {
	data := mypostpage.Data{}
	if u.UserName != "" {
		data.Authenticated = true
		data.Username = u.UserName
	}

	var threadName string
	if r.Method == http.MethodGet {
		threadName = r.URL.Query().Get("name")
	} else if r.Method == http.MethodPost {
		threadName = r.PostFormValue("search")
	} else {
		http.Redirect(w, r, "", http.StatusMethodNotAllowed)
		return
	}

	posts := likedposts.GetPostsByIDs(post.GetPostIDsByTID(post.GetThreadID(threadName)))
	data.Posts = posts
	templates := []string{
		"./web/templates/somemain.html",
		"./web/templates/somenavbar.html",
		"./web/templates/somefooter.html",
	}

	tpl, err := template.ParseFiles(templates...)
	if err != nil {
		log.Println("Failed on SignUp", err.Error())
		funcs.ErrorHandler(w, "Something went wrong 1", http.StatusInternalServerError)
		return
	}

	if err := tpl.Execute(w, data); err != nil {
		funcs.ErrorHandler(w, "Something went wrong 2"+err.Error(), http.StatusInternalServerError)
	}
}
