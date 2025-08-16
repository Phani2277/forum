package models

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

// User represents a registered user of the forum.
type User struct {
	ID       int
	Email    string
	Username string
	Password string
}

// CreateUser inserts a new user into the database with a hashed password.
func CreateUser(db *sql.DB, email, username, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO users(email, username, password) VALUES (?,?,?)", email, username, string(hashed))
	return err
}

// Authenticate checks the credentials and returns the user if valid.
func Authenticate(db *sql.DB, email, password string) (*User, error) {
	u := &User{}
	err := db.QueryRow("SELECT id, email, username, password FROM users WHERE email = ?", email).
		Scan(&u.ID, &u.Email, &u.Username, &u.Password)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, err
	}
	return u, nil
}

// GetByID retrieves a user by ID.
func GetByID(db *sql.DB, id int) (*User, error) {
	u := &User{}
	err := db.QueryRow("SELECT id, email, username FROM users WHERE id = ?", id).
		Scan(&u.ID, &u.Email, &u.Username)
	if err != nil {
		return nil, err
	}
	return u, nil
}
