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

func GetCommentsByPostID(db *sql.DB, postID int) ([]Comment, error) {

	rows, err := db.Query(`
    SELECT id, post_id, user_id, content, created_at, likes, dislikes
    FROM comments
    WHERE post_id = ?
    ORDER BY created_at DESC
`, postID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var comments []Comment

	for rows.Next() {
		var c Comment

		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.Likes, &c.Dislikes)

		if err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}

	return comments, nil
}
