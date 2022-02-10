package post

import (
	"forum/pkg/funcs"
	"forum/pkg/sserver"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"
)

//CreatePost ...
func CreatePost(w http.ResponseWriter, r *http.Request, user *funcs.User) {
	var err error

	Errors := &ErrorsCreatePost{}
	data := &NewData{
		PostData: &PostData{Authenticated: true, CurrentUrl: r.URL.Path, Username: user.UserName},
		Err:      Errors,
	}

	templates := []string{
		"./web/templates/postCreate.html",
		"./web/templates/somenavbar.html",
		"./web/templates/somefooter.html",
	}

	tpl, err := template.ParseFiles(templates...)
	if err != nil {
		log.Println("Failed on CreatePost", err.Error())
		funcs.ErrorHandler(w, "Something went wrong"+err.Error(), 500)
		return
	}

	if r.Method != http.MethodPost {
		tpl.Execute(w, data)
		return
	}

	r.ParseForm()

	thread := Thread{
		Name: r.PostFormValue("thread"),
	}

	post := &Post{
		UserID:       user.ID,
		Author:       user.UserName,
		Title:        r.PostFormValue("title"),
		Content:      r.PostFormValue("postContent"),
		CreationDate: time.Now().Format("January 2 15:04"),
	}

	if strings.ReplaceAll(post.Title, " ", "") == "" {
		Errors.Title = "\"Title\" field is empty"
		tpl.Execute(w, data)
		return
	}

	if strings.ReplaceAll(post.Content, " ", "") == "" {
		Errors.Content = "\"Content\" field is empty"
		tpl.Execute(w, data)
		return
	}

	if thread.Name == "" {
		Errors.Category = "\"Category\" field is empty"
		tpl.Execute(w, data)
		return
	}

	threads := CheckNumberOfThreads(thread.Name)
	if len(threads) == 0 {
		Errors.Category = "\"Category\" must contain only latin letters and arabic numbers"
		tpl.Execute(w, data)
		return
	}

	var errorString string
	post.ID, errorString = InsertPostInfo(post)

	if post.ID == -1 {
		log.Println(errorString)
		funcs.ErrorHandler(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err = NewRating(post.ID)
	if err != nil {
		log.Println(err.Error())
		funcs.ErrorHandler(w, "Something went wrong", http.StatusBadRequest)
		return
	}

	// If post has several threads, to this post will attach this info
	for _, threadName := range threads {
		InsertThreadInfo(threadName, post.ID)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// InsertPostInfo ...
func InsertPostInfo(p *Post) (int, string) {
	stmnt, err := sserver.Db.Db.Prepare("INSERT INTO Posts (userid, author, title, content, creationDate) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return -1, err.Error() + " :Statement Error"
	}
	defer stmnt.Close()

	res, err := stmnt.Exec(p.UserID, p.Author, p.Title, p.Content, p.CreationDate)
	if err != nil {
		return -1, err.Error() + " :Statement executing Error"
	}

	// to get PostID of new post for error handling in DB
	id, err := res.LastInsertId()
	if err != nil {
		return -1, err.Error() + " :Getting ID Error"
	}

	return int(id), ""
}
