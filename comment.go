package main

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID        int
	PostID    int
	UserID    int
	Content   string
	CreatedAt time.Time
	Likes     int
	Dislikes  int
}

func CreateComment(db *sql.DB, postID, userID int, content string) error {
	query := `
        INSERT INTO comments (post_id, user_id, content, created_at)
        VALUES (?, ?, ?, ?)
    `
	_, err := db.Exec(query, postID, userID, content, time.Now())
	return err
}
