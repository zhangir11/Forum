package comment

import (
	"forum/internal/post"
	"forum/pkg/funcs"
	"forum/pkg/sserver"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var ercase = " :Statement Error make comment"

func MakeComment(w http.ResponseWriter, r *http.Request, u *funcs.User) {
	userData := &post.PostData{
		Authenticated: true,
		UserID:        int(u.ID),
		Username:      u.UserName,
		CurrentUrl:    r.URL.Path,
	}
	var err error

	if r.Method != http.MethodPost {
		funcs.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	//recieving comment data
	coment := &post.Commentary{
		AuthorID:     int(u.ID),
		Content:      r.PostFormValue("aComment"),
		CreationDate: time.Now().Format("January 2 15:04"),
	}

	if strings.ReplaceAll(coment.Content, " ", "") == "" {
		funcs.ErrorHandler(w, "Bad request", http.StatusBadRequest)
		return
	}

	if coment.PostID, err = strconv.Atoi(r.URL.Query().Get("id")); err != nil {
		funcs.ErrorHandler(w, "Bad request", http.StatusBadRequest)
		return
	}

	if userData.Post, err = post.GetPostByID(coment.PostID); err != nil {
		funcs.ErrorHandler(w, "Bad request", http.StatusBadRequest)
		return
	}
	// Inserting Comment
	stmnt, err := sserver.Db.Db.Prepare("INSERT INTO Comments (postID, authorID, content, creationDate) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Println(err.Error() + ercase)
		funcs.ErrorHandler(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	res, err := stmnt.Exec(coment.PostID, coment.AuthorID, coment.Content, coment.CreationDate)
	if err != nil {
		log.Println(err.Error() + ercase)
		funcs.ErrorHandler(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	coment.ID = int(id)
	stmnt.Close()

	//Creating Rating
	stmnt, err = sserver.Db.Db.Prepare("INSERT INTO CommentRating (commentID, likeCount, dislikeCount) VALUES (?, 0, 0)")
	if err != nil {
		log.Println(err.Error() + ercase)
		funcs.ErrorHandler(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	_, err = stmnt.Exec(coment.ID)
	if err != nil {
		log.Println(err.Error() + ercase)
		funcs.ErrorHandler(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	//Getting comments
	userData.Comments, _ = post.Get(int(userData.Post.ID))
	post.ShowPost(w, userData)
}
