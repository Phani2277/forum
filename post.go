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

func GetAllPostsd(db *sql.DB) ([]Post, error) {
	query := `SELECT id, user_id, title, content, category, created_at, count_likes FROM posts ORDER BY created_at DESC`

	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.Category, &p.CreatedAt, &p.CountLikes)

		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}
	return posts, nil
}
