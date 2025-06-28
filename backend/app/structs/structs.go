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
	Avatar   sql.NullString
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
	Avatar      Avatar
}

type UsersGet struct {
	ID       ID     `json:"profile_id"`
	Online   bool   `json:"online"`
	Is_Group bool   `json:"is_group"`
	Username string `json:"profile_name"` // Exported field
}

type Group struct {
	UserName Name   `json:"username"`
	Avatar   Avatar `json:"avatar"`
	About    About  `json:"about"`
}

type GroupGet struct {
	ID          ID
	GroupName   Name
	Avatar      sql.NullString
	Description About
}

type GroupReq struct {
	Gid int `json:"gid"`
	Uid int `json:"uid"`
}

type Gusers struct {
	Uid    int
	Name   string
	Avatar sql.NullString
	Adm    bool
}

type Post struct {
	ID           int
	UserId       int
	GroupId      sql.NullInt64
	UserName     string
	GroupName    sql.NullString
	Content      string
	ImagePath    sql.NullString
	CreatedAt    string
	AvatarUser   sql.NullString
	AvatarGroup  sql.NullString
	Privacy      string
	CommentCount int
	LikeCount    int
}

type PostGet struct {
	Start   ID `json:"start"`
	UserId  ID `json:"userId"`
	GroupId ID `json:"grouId"`
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
	Sender     ID        `json:"sender"`
	SenderName string    `json:"author_name"`
	Content    string    `json:"content"`
	Time       time.Time `json:"sent_at"`
}
