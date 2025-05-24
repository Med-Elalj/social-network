package structs

type Post struct {
	ID           int
	Title        string
	Content      string
	CreatedAt    string
	Username     string
	Categories   string
	CommentCount int
	Attachement  string
	Status       string
	// LikeCount    int
	// DislikeCount int
}

type User struct {
	ID        int
	Username  string
	Email     string
	Fname     string
	Lname     string
	Status    string
	Followers []int
	Followed  []int
}
