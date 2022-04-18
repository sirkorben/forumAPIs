package models

import (
	"errors"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
	ErrDuplicateUsername  = errors.New("models: duplicate username")
)

type Post struct {
	Id               int         `json:"-"`
	Title            string      `json:"title"`
	Content          string      `json:"content"`
	CreateDate       string      `json:"creationdate"`
	Categories       []*Category `json:"categories"`
	User             *User       `json:"user"`
	Likes            []*User     `json:"likes"`
	Dislikes         []*User     `json:"dislikes"`
	Comments         []*Comment  `json:"comments"`
	IsLikedByUser    bool        `json:"-"`
	IsDislikedByUser bool        `json:"-"`
}

type User struct {
	Id             int    `json:"-"`
	FirstName      string `json:"firstname"`
	LastName       string `json:"lastname"`
	Age            int    `json:"age"`
	Gender         string `json:"gender"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	HashedPassword []byte `json:"-"`
	Created        string `json:"-"`
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
