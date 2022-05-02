package models

import (
	"errors"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
	ErrDuplicateUsername  = errors.New("models: duplicate username")
	ErrTooManySpaces      = errors.New("inupt data: too many spaces in field")
	InternalServerError   = errors.New("INTERNAL_SERVER_ERROR")
)

type Post struct {
	Id               int         `json:"id,omitempty"`
	Title            string      `json:"title,omitempty"`
	Content          string      `json:"content,omitempty"`
	CreationDate     int         `json:"creation_date,omitempty"`
	Categories       []*Category `json:"categories,omitempty"`
	User             *User       `json:"user,omitempty"`
	Likes            []*User     `json:"likes,omitempty"`
	Dislikes         []*User     `json:"dislikes,omitempty"`
	Comments         []*Comment  `json:"comments,omitempty"`
	IsLikedByUser    bool        `json:"-"`
	IsDislikedByUser bool        `json:"-"`
}

type User struct {
	Id             int    `json:"id,omitempty"`
	FirstName      string `json:"firstname,omitempty"`
	LastName       string `json:"lastname,omitempty"`
	Age            int    `json:"age,omitempty"`
	Gender         string `json:"gender,omitempty"`
	Username       string `json:"username,omitempty"`
	Email          string `json:"email,omitempty"`
	Password       string `json:"password,omitempty"`
	HashedPassword []byte `json:"-"`
	CreationDate   int    `json:"creation_date,omitempty"` // do we need it in user profile?
	// Is_Active      bool		`json:"is_active"`   // for showing user online?
}

type Session struct {
	Id   string
	User *User
}

type Comment struct {
	Id               int
	PostId           int
	User             *User
	Content          string
	IsLikedByUser    bool
	IsDislikedByUser bool
	LikeUsers        []int
	DislikeUsers     []int
	LikeCount        int
	DislikeCount     int
}

type PostReaction struct {
	PostComment     string `json:"postcomment"`
	PostLikeDislike string `json:"postlikedislike"`
}

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Message struct {
	Id           int    `json:"post_id,omitempty"`
	ToFromUser   int    `json:"to_from_user,omitempty"`
	Content      string `json:"content"`
	CreationDate int    `json:"creation_date"`
}
