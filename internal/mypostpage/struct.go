package mypostpage

import "forum/internal/post"

type Data struct {
	Authenticated bool
	Username      string
	Posts         []*post.Post
	CurrentUrl    string
	Page          string
}
