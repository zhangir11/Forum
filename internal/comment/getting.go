package comment

import (
	"forum/internal/post"
	"forum/pkg/sserver"
	"log"
)

//GetByID ...
func GetByID(comentID int) (*post.Commentary, error) {
	pist := &post.Commentary{}
	err := sserver.Db.Db.QueryRow("SELECT * FROM Comments WHERE ID = ?", comentID).Scan(
		&pist.ID,
		&pist.PostID,
		&pist.AuthorID,
		&pist.Content,
		&pist.CreationDate,
	)

	if err != nil {
		log.Println("error on func GetCommentsByPostID() ", err.Error())
		return nil, err
	}

	u, err := sserver.Servak.Database.GetUserByID(pist.AuthorID)
	if err != nil {
		pist.AuthorName = "Deleted"
	} else {
		pist.AuthorName = u.UserName
	}

	pist.Rating = post.GetRatingOfComment(pist.ID)
	return pist, nil
}
