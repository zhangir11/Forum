package mainpage

import (
	"forum/internal/mypostpage"
	"forum/internal/post"
	"forum/pkg/funcs"
	"log"
	"net/http"
	"text/template"
)

//MainPage ...
func MainPage(w http.ResponseWriter, r *http.Request, u *funcs.User) {
	templates := []string{
		"./web/templates/somemain.html",
		"./web/templates/somenavbar.html",
		"./web/templates/somefooter.html",
	}

	tpl, err := template.ParseFiles(templates...)
	if err != nil {
		log.Println("Failed on MainPage", err.Error())
		funcs.ErrorHandler(w, "Something went wrong 1", http.StatusInternalServerError)
		return
	}

	mainPage := &mypostpage.Data{}
	Posts, _ := post.GetPosts()
	mainPage.Posts = Posts

	if u.UserName != "" {
		mainPage.Authenticated = true
		mainPage.Username = u.UserName
	}

	tpl.Execute(w, mainPage)
}
