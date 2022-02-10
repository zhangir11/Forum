package post

import (
	"errors"
	"forum/pkg/sserver"
)

//NewRating ...
func NewRating(id int) error {
	stmnt, err := sserver.Db.Db.Prepare("INSERT INTO PostRating (postID, likeCount, dislikeCount) VALUES (?, ?, ?)")
	if err != nil {
		return errors.New(err.Error() + " :Statement Error")
	}
	defer stmnt.Close()

	_, err = stmnt.Exec(id, 0, 0)
	if err != nil {
		return errors.New(err.Error() + " :Statement executing Error")
	}

	return nil
}

// GetRateCountOfPost ...
func GetRateCountOfPost(postID int) *Rating {
	rates := NewPostRating()
	if err := sserver.Db.Db.QueryRow("SELECT * FROM PostRating WHERE postID = ?", postID).
		Scan(&rates.PostID,
			&rates.LikeCount,
			&rates.DislikeCount,
		); err != nil {
		// It means nobody rated the post, likeCount and dislikeCount now is zero
		rates.LikeCount = 0
		rates.DislikeCount = 0
	}
	return rates
}

//GetRatingOfComment ...
func GetRatingOfComment(commentID int) *CommentRating {
	res := &CommentRating{}
	sserver.Db.Db.QueryRow("SELECT likeCount, dislikeCount FROM CommentRating WHERE commentID = ?", commentID).Scan(
		&res.LikeCount,
		&res.DislikeCount,
	)
	return res
}
