package domain

type Comment struct {
	ID       string
	PostID   string
	ParentID *string
	Author   string
	Text     string
	Children []*Comment
}

type Post struct {
	ID       string
	Title    string
	Content  string
	Author   string
	Flag     bool
	Comments []*Comment
}
