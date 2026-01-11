package repo

import (
	"database/sql"
	"time"

	"forum/internal/models"
)

func CreateComment(db *sql.DB, postID int, userID int, content string) error {
	query := `
        INSERT INTO comments (post_id, user_id, content, created_at)
        VALUES (?, ?, ?, ?)
    `
	_, err := db.Exec(query, postID, userID, content, time.Now())
	return err
}

func GetCommentsByPostID(db *sql.DB, postID int) ([]models.CommentView, error) {
	query := `
    SELECT c.id, c.user_id, u.username, c.content, c.created_at
    FROM comments c
    JOIN users u ON u.id = c.user_id
    WHERE c.post_id = ?
    ORDER BY c.created_at DESC
`
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.CommentView
	for rows.Next() {
		var c models.CommentView
		if err := rows.Scan(&c.ID, &c.UserID, &c.AuthorName, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func CommentExists(db *sql.DB, commentID int) (bool, error) {
	row := db.QueryRow(`SELECT 1 FROM comments WHERE id = ? LIMIT 1`, commentID)
	var one int
	err := row.Scan(&one)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
