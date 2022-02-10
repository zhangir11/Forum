package post

// Post ...
type Post struct {
	ID           int    `db:"postid"`
	UserID       int    `db:"userid"`
	Title        string `db:"title"`
	Author       string `db:"author"`
	Content      string `db:"content"`
	CreationDate string `db:"creationDate"`
	Threads      []*Thread
	PostRate     *Rating
}

// NewPost ...
func NewPost() *Post {
	return &Post{}
}

// Thread ...
type Thread struct {
	ID   int64  `db:"ID"`
	Name string `db:"Name"`
}

// NewThread ....
func NewThread() *Thread {
	return &Thread{}
}

// Rating ...
type Rating struct {
	PostID       int64 `db:"postID"`
	LikeCount    int64 `db:"likeCount"`
	DislikeCount int64 `db:"dislikeCount"`
}

// NewPostRating ...
func NewPostRating() *Rating {
	return &Rating{}
}

//ErrorsCreatePost ...
type ErrorsCreatePost struct {
	Title, Content, Category, InvalidInput string
}

type NewData struct {
	*PostData
	Err *ErrorsCreatePost
}

// Commentary is comment
type Commentary struct {
	AuthorName   string
	ID           int    `db:"ID"`
	PostID       int    `db:"postID"`
	AuthorID     int    `db:"authorID"`
	Content      string `db:"content"`
	CreationDate string `db:"creationDate"`
	Rating       *CommentRating
}

// Rating ...
type CommentRating struct {
	CommentID    int `db:"commentID"`
	PostID       int `db:"postID"`
	LikeCount    int `db:"likeCount"`
	DislikeCount int `db:"dislikeCount"`
}

// PostData : stores the data for template execution
type PostData struct {
	Authenticated bool
	Post          *Post
	Comments      []*Commentary
	UserID        int
	CurrentUrl    string
	Username      string
}
