package repo

import (
	"database/sql"

	"forum/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(db *sql.DB, email string, username string, password string) error {
	query := `INSERT INTO users (email, username, password) VALUES (?, ?, ?)`

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec(query, email, username, string(hash))
	return err
}

func GetUserByEmail(db *sql.DB, email string) (*models.User, error) {
	query := `SELECT id, email, username, password FROM users WHERE email = ? LIMIT 1`

	row := db.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByID(db *sql.DB, id int) (*models.User, error) {
	query := `SELECT id, email, username, password FROM users WHERE id = ? LIMIT 1`

	row := db.QueryRow(query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
