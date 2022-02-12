package baselogic

import (
	"errors"
	"fmt"
	"forum/internal/comment"
	"forum/internal/delrate"
	"forum/internal/filterbythread"
	"forum/internal/likedposts"
	"forum/internal/mainpage"
	"forum/internal/mypostpage"
	"forum/internal/post"
	"forum/internal/signupin"
	"forum/internal/testingauthgoogle"
	"forum/pkg/sserver"
	"net/http"
)

var (
	AuthIsRequired   = 1
	AuthIsIrrelevant = 0
	AuthIsForbidden  = 2
)

// Start - Initializing server
func Start(s *sserver.Server) error {
	testingauthgoogle.Init()
	sserver.DeclareServak(s)
	if s == nil {
		return errors.New("error: server is nil")
	}
	multiHandlerer := http.NewServeMux()
	fmt.Println("Server is working on port " + s.WebPort)
	fmt.Println("http://localhost:" + s.WebPort + "/")
	multiHandlerer.Handle("/testgoogle", smth(testingauthgoogle.HandleGoogleLogin))
	multiHandlerer.Handle("/callback", smth(testingauthgoogle.HandleGoogleCallback))
	multiHandlerer.Handle("/likedposts", Middleware(likedposts.ShowLovely, AuthIsRequired))
	multiHandlerer.Handle("/tags", Middleware(filterbythread.FilterByThread, AuthIsIrrelevant))
	multiHandlerer.Handle("/myposts", Middleware(mypostpage.MyPostPage, AuthIsRequired))
	multiHandlerer.Handle("/main", Middleware(mainpage.MainPage, AuthIsIrrelevant))
	multiHandlerer.Handle("/createcomment", Middleware(comment.MakeComment, AuthIsRequired))
	multiHandlerer.Handle("/post", Middleware(post.Postpage, AuthIsIrrelevant))
	multiHandlerer.Handle("/signup", Middleware(signupin.SignUp, AuthIsForbidden))
	multiHandlerer.Handle("/delete", Middleware(delrate.DeleteSmth, AuthIsRequired))
	multiHandlerer.Handle("/rate", Middleware(delrate.RateSmth, AuthIsRequired))
	multiHandlerer.Handle("/signin", Middleware(signupin.Login, AuthIsForbidden))
	multiHandlerer.Handle("/create", Middleware(post.CreatePost, AuthIsRequired))
	multiHandlerer.Handle("/logout", Middleware(Logout, AuthIsRequired))
	multiHandlerer.Handle("/login", Middleware(signupin.Login, AuthIsForbidden))
	multiHandlerer.Handle("/", Middleware(mainpage.MainPage, AuthIsIrrelevant))

	return http.ListenAndServe(":"+s.WebPort, multiHandlerer)

}

func smth(f func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r)
	}
}
