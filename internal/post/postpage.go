package post

import (
	"fmt"
	"forum/pkg/funcs"
	"forum/pkg/sserver"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

func Postpage(w http.ResponseWriter, r *http.Request, u *funcs.User) {
	strID := r.URL.Query().Get("id")

	id, err := strconv.Atoi(strID)
	if err != nil {
		funcs.ErrorHandler(w, "something went wrong 2", http.StatusBadRequest)
		return
	}

	post, notFound := GetPostByID(id)

	if notFound != nil {
		funcs.ErrorHandler(w, "Something went wrong 3", http.StatusNotFound)
		return
	}

	Comments, _ := Get(id)

	userData := &PostData{
		Post:       post,
		Comments:   Comments,
		CurrentUrl: r.URL.Path,
	}

	fmt.Println("postID", id)

	if u.UserName != "" {
		userData.Authenticated = true
		userData.Username = u.UserName

		if r.Method == http.MethodPost {
			str := r.PostFormValue("submit")
			kind := 1

			if str == "dislike" {
				kind = -1
			}
			RatePostFunc(post, kind, u)
		}
	}

	ShowPost(w, userData)
}

func ShowPost(w http.ResponseWriter, userData *PostData) {
	templates := []string{
		"./web/templates/post.html",
		"./web/templates/somenavbar.html",
		"./web/templates/somefooter.html",
	}

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		funcs.ErrorHandler(w, "something went wrong 1"+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, userData); err != nil {
		funcs.ErrorHandler(w, "Something went wrong 5", http.StatusInternalServerError)
		log.Println(err.Error())
	}
}

//Get Comments by PostID
func Get(postID int) ([]*Commentary, error) {
	pists := []*Commentary{}

	rows, err := sserver.Db.Db.Query("SELECT * FROM Comments WHERE postID = ?", postID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		pist := &Commentary{PostID: postID}

		if err = rows.Scan(
			&pist.ID,
			&pist.PostID,
			&pist.AuthorID,
			&pist.Content,
			&pist.CreationDate,
		); err != nil {
			log.Println("error on func GetCommentsByPostID() ", err.Error())
			return nil, err
		}

		u, err := sserver.Servak.Database.GetUserByID(pist.AuthorID)
		if err != nil {
			pist.AuthorName = "Deleted"
		} else {
			pist.AuthorName = u.UserName
		}

		pist.Rating = GetRatingOfComment(pist.ID)
		pists = append(pists, pist)
	}
	return pists, nil
}
