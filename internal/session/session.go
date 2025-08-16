package session

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

// Duration for session validity.
const Duration = 24 * time.Hour

// Create creates a new session for the given user and returns the token.
func Create(db *sql.DB, userID int) (string, error) {
	token := uuid.Must(uuid.NewV4()).String()
	expires := time.Now().Add(Duration)
	_, err := db.Exec("INSERT INTO sessions(user_id, token, expires_at) VALUES (?,?,?)", userID, token, expires)
	if err != nil {
		return "", err
	}
	return token, nil
}

// GetUserID returns the user ID associated with the token. If the session is
// expired it is removed and an error is returned.
func GetUserID(db *sql.DB, token string) (int, error) {
	var userID int
	var expires time.Time
	err := db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE token = ?", token).Scan(&userID, &expires)
	if err != nil {
		return 0, err
	}
	if time.Now().After(expires) {
		db.Exec("DELETE FROM sessions WHERE token = ?", token)
		return 0, errors.New("session expired")
	}
	return userID, nil
}

// Delete removes a session token from the database.
func Delete(db *sql.DB, token string) {
	db.Exec("DELETE FROM sessions WHERE token = ?", token)
}
