package structs

import "time"

type PostInfo struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
}

type CommentInfo struct {
	PostID  string `json:"post_id"`
	Content string `json:"content"`
}

type Post1 struct {
	Pid          int      `json:"pid"`
	Title        string   `json:"title"`
	Content      string   `json:"content"`
	Author       string   `json:"author"`
	CreationTime string   `json:"creation_time"`
	Categories   []string `json:"categories"`
}

type Comment struct {
	Pid          int    `json:"pid"`
	Author       string `json:"author"`
	Content      string `json:"content"`
	CreationTime string `json:"creation_time"`
}

type User1 struct {
	Online   bool   `json:"online"`
	Username string `json:"username"` // Exported field
}

type Message struct {
	Sender  string    `json:"sender"`
	Content string    `json:"message"`
	Time    time.Time `json:"time"`
}
