package models

import (
	"database/sql"
	"testing"

	"forum/internal/database"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := database.Init(db); err != nil {
		t.Fatalf("init: %v", err)
	}
	return db
}

func TestCreateAndAuthenticateUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	if err := CreateUser(db, "a@example.com", "alice", "secret"); err != nil {
		t.Fatalf("create user: %v", err)
	}
	user, err := Authenticate(db, "a@example.com", "secret")
	if err != nil {
		t.Fatalf("auth: %v", err)
	}
	if user.Username != "alice" {
		t.Fatalf("unexpected username %s", user.Username)
	}
	if _, err := Authenticate(db, "a@example.com", "wrong"); err == nil {
		t.Fatalf("expected auth failure")
	}
}
