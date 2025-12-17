package main

import (
	"database/sql"
	"time"
)

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

func CreatePost(db *sql.DB, user_id int, title, content, categoryID int) error {
	query := `INSERT INTO posts (user_id, title, content, category_id, created_at) VALUES (?, ?, ?, ?, ?)`

	_, err := db.Exec(query, user_id, title, content, categoryID, time.Now())

	return err
}

func GetAllPosts(db *sql.DB) ([]PostView, error) {
	query := `
    SELECT 
        p.id, p.user_id, p.title, p.content,
        p.category_id, c.name,
        p.created_at, p.count_likes
    FROM posts p
    JOIN categories c ON c.id = p.category_id
    ORDER BY p.created_at DESC
    `
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []PostView
	for rows.Next() {
		var p PostView
		err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CategoryID, &p.CategoryName, &p.CreatedAt, &p.CountLikes)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func GetPostReactionCounts(db *sql.DB, postID int) (int, int, error) {
	row := db.QueryRow(`SELECT
  COALESCE(SUM(CASE WHEN value = 1 THEN 1 ELSE 0 END), 0),
  COALESCE(SUM(CASE WHEN value = -1 THEN 1 ELSE 0 END), 0)
FROM post_reactions
WHERE post_id = ?;`, postID)

	var likes int
	var dislikes int
	err := row.Scan(&likes, &dislikes)

	if err != nil {
		return 0, 0, err
	}

	return likes, dislikes, nil

}
