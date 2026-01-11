package models

import "time"

type CommentView struct {
	ID         int
	UserID     int
	AuthorName string
	Content    string
	CreatedAt  time.Time
	Likes      int
	Dislikes   int
}
