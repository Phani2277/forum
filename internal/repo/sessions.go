package repo

import (
	"database/sql"
	"fmt"
	"time"

	"forum/internal/models"
	"github.com/google/uuid"
)

func CreateSessions(db *sql.DB, userID int) (string, error) {
	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(20 * time.Minute)

	_, err := db.Exec(`DELETE FROM sessions WHERE user_id = ?`, userID)
	if err != nil {
		return "", err
	}

	query := `INSERT INTO sessions(id, user_id, expires_at) VALUES (?, ?, ?)`
	_, err = db.Exec(query, sessionID, userID, expiresAt)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func DeleteSession(db *sql.DB, sessionID string) error {
	_, err := db.Exec(`DELETE FROM sessions WHERE id = ?`, sessionID)
	return err
}

func GetUserBySessionID(db *sql.DB, sessionID string) (*models.User, error) {
	query := `SELECT user_id, expires_at FROM sessions WHERE id = ? LIMIT 1`
	row := db.QueryRow(query, sessionID)

	var userID int
	var expiresAt time.Time

	err := row.Scan(&userID, &expiresAt)
	if err != nil {
		return nil, err
	}

	if time.Now().After(expiresAt) {
		return nil, fmt.Errorf("session expired")
	}

	userQuery := `SELECT id, email, username, password FROM users WHERE id = ? LIMIT 1`
	userRow := db.QueryRow(userQuery, userID)

	var user models.User
	err = userRow.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
