package middleware

import (
	"database/sql"
	"net/http"

	"forum/internal/models"
	"forum/internal/repo"
)

func CurrentUser(db *sql.DB, r *http.Request) (*models.User, error) {
	c, err := r.Cookie("session")
	if err != nil {
		return nil, err
	}
	return repo.GetUserBySessionID(db, c.Value)
}
