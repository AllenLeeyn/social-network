package dbTools

import (
	"database/sql"
	"time"
)

type User struct {
	ID            int
	TypeID        int
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	NickName      string `json:"nickName"`
	Gender        string `json:"gender"`
	Age           int    `json:"age"`
	Email         string `json:"email"`
	Passwd        string `json:"password"`
	ConfirmPasswd string `json:"confirmPassword"`
	PwHash        []byte
	RegDate       time.Time
	LastLogin     time.Time
}

type Session struct {
	ID         string
	UserID     int
	IsActive   bool
	StartTime  time.Time
	ExpireTime time.Time
	LastAccess time.Time
}

type Post struct {
	ID           int       `json:"ID"`
	UserID       int       `json:"userID"`
	UserName     string    `json:"userName"`
	CommentCount int       `json:"commentCount"`
	LikeCount    int       `json:"likeCount"`
	DislikeCount int       `json:"dislikeCount"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"createdAt"`
	Categories   []int     `json:"categories"`
	CatNames     string    `json:"catNames"`
	Rating       int       `json:"rating"`
}

type Comment struct {
	ID           int
	UserID       int    `json:"userID"`
	UserName     string `json:"userName"`
	ParentID     sql.NullInt64
	PostID       int       `json:"postID"`
	Content      string    `json:"content"`
	LikeCount    int       `json:"likeCount"`
	DislikeCount int       `json:"dislikeCount"`
	CreatedAt    time.Time `json:"createdAt"`
	Rating       int       `json:"rating"`
}

type Feedback struct {
	Tgt       string `json:"tgt"`
	UserID    int
	ParentID  int `json:"parentID"`
	Rating    int `json:"rating"`
	CreatedAt time.Time
}

type Message struct {
	ID         int          `json:"ID"`
	Action     string       `json:"action"`
	SenderID   int          `json:"senderID"`
	ReceiverID int          `json:"receiverID"`
	Content    string       `json:"content"`
	CreatedAt  time.Time    `json:"createdAt"`
	ReadAt     sql.NullTime `json:"readAt"`
}
