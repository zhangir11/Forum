package likedposts

import (
	"forum/internal/mypostpage"
	"forum/internal/post"
	"forum/pkg/funcs"
	"forum/pkg/sserver"
	"log"
	"net/http"
	"text/template"
)

func ShowLovely(w http.ResponseWriter, r *http.Request, u *funcs.User) {
	templates := []string{
		"./web/templates/somemain.html",
		"./web/templates/somenavbar.html",
		"./web/templates/somefooter.html",
	}

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		log.Println(err)
		funcs.ErrorHandler(w, "Something Went Wrong", http.StatusInternalServerError)
		return
	}

	posts := GetPostsByIDs(GetIDLovelyPosts(u.ID))
	data := &mypostpage.Data{Posts: posts, Authenticated: true, CurrentUrl: r.URL.Path, Page: "Liked Posts", Username: u.UserName}

	if err = tmpl.Execute(w, data); err != nil {
		funcs.ErrorHandler(w, "Something Went Wrong", http.StatusInternalServerError)
		return
	}
}

func GetPostsByIDs(IDs []int) []*post.Post {
	res := make([]*post.Post, len(IDs))
	for ind, ID := range IDs {
		res[ind], _ = post.GetPostByID(ID)
	}
	return res
}

func GetIDLovelyPosts(AuthorID int) []int {
	result := []int{}
	// show posts in reverse order (in the beginning fresh posts)
	res, err := sserver.Db.Db.Query("SELECT postID FROM RateUserPost WHERE kind = 1 AND userID = ?", AuthorID)
	if err != nil {
		return nil
	}

	defer res.Close()

	for res.Next() {
		postID := 0
		if err := res.Scan(&postID); err != nil {
			log.Println(err.Error(), "GetIDLovelyPosts() postsStore.go")
			return nil
		}
		result = append(result, postID)
	}

	return result
}
