package post

import (
	"database/sql"
	"forum/pkg/funcs"
	"forum/pkg/sserver"
	"log"
)

//GetPostsByAuthor ...
func GetPostsByAuthor(u *funcs.User) []*Post {
	// show posts in reverse order (in the beginning fresh posts)
	res, err := sserver.Db.Db.Query("SELECT * FROM Posts WHERE userid = ?", u.ID) // DESC - in reverse order
	if err != nil {
		return nil
	}
	defer res.Close()

	return getpostsbySQLresult(res)
}

func getpostsbySQLresult(res *sql.Rows) []*Post {
	var posts []*Post
	for res.Next() {
		post1 := NewPost()
		if err := res.Scan(&post1.ID, &post1.UserID, &post1.Author, &post1.Title, &post1.Content, &post1.CreationDate); err != nil {
			log.Println(err.Error(), "GetPosts() postsStore.go")
			return nil

		}
		post1.Threads, _ = GetThreadOfPost(post1.ID)
		post1.PostRate = GetRateCountOfPost(post1.ID)
		posts = append(posts, post1)
	}
	return posts
}

// GetPosts ...
func GetPosts() ([]*Post, error) {
	// show posts in reverse order (in the beginning fresh posts)
	res, err := sserver.Db.Db.Query("SELECT * FROM Posts ORDER BY postid DESC") // DESC - in reverse order
	if err != nil {
		return nil, err
	}
	defer res.Close()

	return getpostsbySQLresult(res), nil
}

// GetThreadByID ...
func GetThreadByID(id int) (*Thread, error) {
	thread := NewThread()
	if err := sserver.Db.Db.QueryRow("SELECT * FROM Threads WHERE ID = ?", id).Scan(&thread.ID, &thread.Name); err != nil {
		log.Println("error on func GetThreadByID()")
		return nil, err
	}
	return thread, nil
}

//GetPostByID ...
func GetPostByID(id int) (*Post, error) {
	pist := &Post{}
	if err := sserver.Db.Db.QueryRow("SELECT * FROM Posts WHERE postid = ?", id).Scan(
		&pist.ID,
		&pist.UserID,
		&pist.Author,
		&pist.Title,
		&pist.Content,
		&pist.CreationDate,
	); err != nil {
		log.Println("error on func GetPostByID() ", err.Error())
		return nil, err
	}
	pist.Threads, _ = GetThreadOfPost(pist.ID)
	pist.PostRate = GetRateCountOfPost(pist.ID)
	return pist, nil
}
