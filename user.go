package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       int
	Email    string
	Username string
	Password string
}

func CreateUser(db *sql.DB, email string, username string, password string) error {
	query := `INSERT INTO users (email, username, password) VALUES (?, ?, ?)`

	_, err := db.Exec(query, email, username, password)

	return err
}

func CreateSessions(db *sql.DB, userID int) (string, error) {
	sessionID := uuid.New().String()

	expiresAt := time.Now().Add(20 * time.Minute)

	query := `INSERT INTO sessions(id, user_id, expires_at) VALUES (?, ?, ?)`
	_, err := db.Exec(query, sessionID, userID, expiresAt)

	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	query := `SELECT id, email, username, password FROM users WHERE email = ? LIMIT 1`

	row := db.QueryRow(query, email)

	var user User

	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserBySessionID(db *sql.DB, sessionID string) (*User, error) {
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

	var user User

	err = userRow.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByID(db *sql.DB, id int) (*User, error) {
	query := `SELECT id, email, username, password FROM users WHERE id = ? LIMIT 1`

	row := db.QueryRow(query, id)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
