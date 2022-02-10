package delrate

import (
	"database/sql"
	"errors"
	"forum/internal/comment"
	"forum/internal/post"
	"forum/pkg/funcs"
	"forum/pkg/sserver"
	"net/http"
	"strconv"
)

//RateSmth ...
func RateSmth(w http.ResponseWriter, r *http.Request, u *funcs.User) {
	tempID := r.URL.Query().Get("coment_id")

	ID, err := strconv.Atoi(tempID)
	if err == nil {
		if RateComment(w, r, ID, u) != nil {
			funcs.ErrorHandler(w, "Something went wrong", http.StatusInternalServerError)
		}
		return
	}

	tempID = r.URL.Query().Get("post_id")

	ID, err = strconv.Atoi(tempID)
	if err != nil {
		funcs.ErrorHandler(w, "BadRequest", http.StatusBadRequest)
		return
	}

	if RatePost(w, r, ID, u) != nil {
		funcs.ErrorHandler(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}

//RateComment ...
func RateComment(w http.ResponseWriter, r *http.Request, ID int, u *funcs.User) error {
	tempkind := 1

	if ID < 0 {
		ID = -ID
		tempkind = -1
	}

	coment, err := comment.GetByID(ID)

	if err != nil {
		return err
	}

	var kind int
	var s1, s2 *sql.Stmt

	err = sserver.Db.Db.QueryRow("SELECT kind FROM RateUserComment WHERE commentID = ? AND userID = ?", ID, u.ID).Scan(&kind)
	if err != nil {
		s1, err = sserver.Db.Db.Prepare("INSERT INTO RateUserComment (kind, commentID, userID) VALUES (?, ?, ?)")

		if err != nil {
			return err
		}

		if tempkind < 0 {
			coment.Rating.DislikeCount++
		} else {
			coment.Rating.LikeCount++
		}

	} else if kind == tempkind {
		http.Redirect(w, r, "/post?id="+strconv.Itoa(int(coment.PostID)), http.StatusSeeOther)
		return nil
	} else {

		s1, err = sserver.Db.Db.Prepare("UPDATE RateUserComment SET kind = ? WHERE commentID= ? AND userID = ?")
		if err != nil {
			return err
		}

		if tempkind < 0 {
			coment.Rating.LikeCount--
			coment.Rating.DislikeCount++
		} else {
			coment.Rating.LikeCount++
			coment.Rating.DislikeCount--
		}
	}

	_, err = s1.Exec(tempkind, ID, u.ID)
	if err != nil {
		return err
	}
	s1.Close()

	s2, err = sserver.Db.Db.Prepare("UPDATE CommentRating SET likeCount = ?, dislikeCount = ? WHERE commentID= ?")
	if err != nil {
		return err
	}

	_, err = s2.Exec(coment.Rating.LikeCount, coment.Rating.DislikeCount, ID)
	if err != nil {
		return err
	}
	s2.Close()

	http.Redirect(w, r, "/post?id="+strconv.Itoa(int(coment.PostID)), http.StatusSeeOther)
	return nil
}

//RatePost ...
func RatePost(w http.ResponseWriter, r *http.Request, ID int, u *funcs.User) error {
	kind := 1
	if ID < 0 {
		kind = -1
		ID = -ID
	}

	posta, err := post.GetPostByID(ID)
	if err != nil {
		return errors.New(err.Error() + " ratePost")
	}

	err = post.RatePostFunc(posta, kind, u)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return err
}
