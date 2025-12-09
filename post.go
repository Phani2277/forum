package main

import (
	"database/sql"
	"time"
)

type Post struct {
	ID         int
	UserID     int
	Title      string
	Content    string
	Category   string
	CreatedAt  time.Time
	CountLikes int
}

func CreatePost(db *sql.DB, user_id int, title, content, category string) error {
	query := `INSERT INTO posts (user_id, title, content, category, created_at) VALUES (?, ?, ?, ?, ?)`

	_, err := db.Exec(query, user_id, title, content, category, time.Now())

	return err
}
