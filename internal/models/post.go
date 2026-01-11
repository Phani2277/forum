package models

import "time"

type CommentCard struct {
	AuthorName string
	Content    string
}

type BasePageData struct {
	CurrentUser *User
	Error       string
}

type PostView struct {
	ID           int
	UserID       int
	Title        string
	Content      string
	CategoryID   int
	CategoryName string
	CreatedAt    time.Time
	CountLikes   int
}

type PostCard struct {
	ID           int
	Title        string
	Content      string
	CategoryName string
	AuthorName   string
	Likes        int
	Dislikes     int
	Comments     []CommentCard
}

type HomePageData struct {
	CurrentUser       *User
	Categories        []Category
	Posts             []PostCard
	AllActive         bool
	MineActive        bool
	LikedActive       bool
	SelectedCategoryID int
}

type PostCardWithComments struct {
	ID           int
	Title        string
	Content      string
	CategoryName string
	AuthorName   string
	Likes        int
	Dislikes     int
	Comments     []CommentView
}

type PostPageData struct {
	CurrentUser *User
	Post        PostCardWithComments
}

type CreatePostPageData struct {
	CurrentUser *User
	Categories  []Category
	Error       string
}
