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

	Avatar sql.NullString
)

type Login struct {
	NoE      string `json:"login"`
	Password string `json:"pwd"`
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
	GroupName string         `json:"groupName"`
	Avatar    sql.NullString `json:"avatar"`
	About     string         `json:"about"`
	Privacy   PostPrivacy    `json:"privacy"`
}

type GroupGet struct {
	ID          ID
	GroupName   string
	Avatar      sql.NullString
	Description string
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
	IsLiked      bool
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
	ImagePath   PImage         `json:"image_path"`
}

type CommentGet struct {
	Pid          ID        `json:"pid"`
	Author       string    `json:"author"`
	Content      string    `json:"content"`
	CreationTime time.Time `json:"creation_time"`
}

type LikeInfo struct {
	EntityID   ID     `json:"entity_id"`
	EntityType string `json:"entity_type"`
	IsLiked    bool   `json:"is_liked"`
}

type Message struct {
	Sender     ID        `json:"sender"`
	SenderName string    `json:"author_name"`
	Content    string    `json:"content"`
	Time       time.Time `json:"sent_at"`
}
