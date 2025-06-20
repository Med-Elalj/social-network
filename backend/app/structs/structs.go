package structs

import (
	"database/sql"
	"time"
)

type (
	PImage      string
	Pbody       string
	PostPrivacy string

	NameOrEmail string

	ID             int
	CommentContent string

	Name     string
	Email    string
	Password string
	Gender   string
	Avatar   string
	About    string
)

type Register struct {
	UserName  Name     `json:"username"`
	Email     Email    `json:"email"`
	Birthdate string   `json:"birthdate"`
	Fname     Name     `json:"fname"`
	Lname     Name     `json:"lname"`
	Password  Password `json:"password"`
	Gender    Gender   `json:"gender"`
	Avatar    Avatar   `json:"avatar"`
	About     About    `json:"about"`
}

type Login struct {
	NoE      NameOrEmail `json:"login"`
	Password Password    `json:"pwd"`
}

type User struct {
	ID          int
	Username    string
	Email       string
	Fname       string
	Lname       string
	Status      string
	Followers   []int
	Followed    []int
	Gender      int
	Description sql.NullString
	IsPublic    bool
	IsPerson    bool
	Avatar      sql.NullString
}

type UsersGet struct {
	Online   bool   `json:"online"`
	Username string `json:"username"` // Exported field
}

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
	LikeCount    int
	DislikeCount int
}

type PostGet struct {
	Pid          ID          `json:"pid"`
	AuthorId     ID          `json:"authorId"`
	Author       string      `json:"author_username"`
	GroupId      ID          `json:"groupId"`
	GroupName    string      `json:"group_name"`
	Title        string      `json:"title"`
	Content      string      `json:"content"`
	CreationTime time.Time   `json:"creation_time"`
	Privacy      PostPrivacy `json:"privacy"`
}

type PostCreate struct {
	Content Pbody       `json:"content"`
	Image   PImage      `json:"image"`
	Privacy PostPrivacy `json:"privacy"`
}

type CommentInfo struct {
	PostID  ID             `json:"post_id"`
	Content CommentContent `json:"content"`
}

type CommentGet struct {
	Pid          ID        `json:"pid"`
	Author       string    `json:"author"`
	Content      string    `json:"content"`
	CreationTime time.Time `json:"creation_time"`
}

type Message struct {
	Sender  ID        `json:"sender"`
	Content string    `json:"message"`
	Time    time.Time `json:"time"`
}
