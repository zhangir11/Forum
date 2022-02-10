package post

import (
	"database/sql"
	"forum/pkg/funcs"
	"forum/pkg/sserver"
)

//RatePostFunc ...
func RatePostFunc(posst *Post, tempkind int, u *funcs.User) error {
	var s1, s2 *sql.Stmt
	var kind int

	err := sserver.Db.Db.QueryRow("SELECT kind FROM RateUserPost WHERE postID = ? AND userID = ?", posst.ID, u.ID).Scan(&kind)
	if err != nil {

		s1, err = sserver.Db.Db.Prepare("INSERT INTO RateUserPost (kind, postID, userID) VALUES (?, ?, ?)")
		if err != nil {
			return err
		}

		s2, err = sserver.Db.Db.Prepare("UPDATE PostRating SET likeCount = ?, dislikeCount = ? WHERE postID= ? ")
		if err != nil {
			return err
		}

		if tempkind < 0 {
			posst.PostRate.DislikeCount++
		} else {
			posst.PostRate.LikeCount++
		}
	} else if kind == tempkind {
		return nil
	} else {
		s1, err = sserver.Db.Db.Prepare("UPDATE RateUserPost SET kind = ? WHERE postID= ? AND userID = ?")
		if err != nil {
			return err
		}

		if tempkind < 0 {
			posst.PostRate.LikeCount--
			posst.PostRate.DislikeCount++
		} else {
			posst.PostRate.LikeCount++
			posst.PostRate.DislikeCount--
		}

		s2, err = sserver.Db.Db.Prepare("UPDATE PostRating SET likeCount = ?, dislikeCount = ? WHERE postID= ? ")
		if err != nil {
			return err
		}

	}

	_, err = s1.Exec(tempkind, posst.ID, u.ID)
	if err != nil {
		return err
	}
	s1.Close()

	_, err = s2.Exec(posst.PostRate.LikeCount, posst.PostRate.DislikeCount, posst.ID)
	if err != nil {
		return err
	}
	s2.Close()
	return nil
}
