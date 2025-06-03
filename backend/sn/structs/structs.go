package structs

import (
	"database/sql"
	"time"
)

type Input interface {
	IsValid() error
}

type NameOrEmail struct {
	Input
}

// Custom field types implementing Validator

type (
	Ptitle      string
	Pbody       string
	Pcategories []string
	PostPrivacy int8

	ID             int
	CommentContent string

	Name      string
	Email     string
	Password  string
	Birthdate time.Time
	Gender    int
	Avatar    string
	About     string
)

// User struct with custom types
type Register struct {
	UserName  Name      `json:"username"`
	Email     Email     `json:"email"`
	Birthdate Birthdate `json:"birthdate"`
	Fname     Name      `json:"fname"`
	Lname     Name      `json:"lname"`
	Password  Password  `json:"password"`
	Gender    Gender    `json:"gender"`
	Avatar    Avatar    `json:"avatar"`
	About     About     `json:"about"`
}

type Login struct {
	NoE      NameOrEmail `json:"login"`
	Password Password    `json:"pwd"`
}

// Input interface with IsValid method

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

type PostCreate struct {
	Title      Ptitle      `json:"title"`
	Content    Pbody       `json:"content"`
	Categories Pcategories `json:"-"`
}

type CommentInfo struct {
	PostID  ID             `json:"post_id"`
	Content CommentContent `json:"content"`
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
	// Categories   []string `json:"categories"` // TODO implement categories
}

type CommentGet struct {
	Pid          ID        `json:"pid"`
	Author       string    `json:"author"`
	Content      string    `json:"content"`
	CreationTime time.Time `json:"creation_time"`
}

type UsersGet struct {
	Online   bool   `json:"online"`
	Username string `json:"username"` // Exported field
}

type Message struct {
	Sender  ID        `json:"sender"`
	Content string    `json:"message"`
	Time    time.Time `json:"time"`
}
