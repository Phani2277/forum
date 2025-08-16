package models

import (
	"database/sql"
	"time"
)

// Comment represents a comment on a post.
type Comment struct {
	ID        int
	PostID    int
	UserID    int
	Username  string
	Content   string
	CreatedAt time.Time
	Likes     int
	Dislikes  int
}

// CreateComment inserts a new comment.
func CreateComment(db *sql.DB, postID, userID int, content string) error {
	_, err := db.Exec("INSERT INTO comments(post_id, user_id, content) VALUES (?,?,?)", postID, userID, content)
	return err
}

// GetComments returns comments for a given post.
func GetComments(db *sql.DB, postID int) ([]Comment, error) {
	rows, err := db.Query(`SELECT c.id, c.post_id, c.user_id, u.username, c.content, c.created_at,
        IFNULL(SUM(CASE WHEN cl.value=1 THEN 1 END),0) as likes,
        IFNULL(SUM(CASE WHEN cl.value=-1 THEN 1 END),0) as dislikes
        FROM comments c
        JOIN users u ON c.user_id=u.id
        LEFT JOIN comment_likes cl ON c.id=cl.comment_id
        WHERE c.post_id=? GROUP BY c.id ORDER BY c.created_at ASC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []Comment
	for rows.Next() {
		cm := Comment{}
		if err := rows.Scan(&cm.ID, &cm.PostID, &cm.UserID, &cm.Username, &cm.Content, &cm.CreatedAt, &cm.Likes, &cm.Dislikes); err != nil {
			return nil, err
		}
		comments = append(comments, cm)
	}
	return comments, rows.Err()
}
