package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		panic("База не открылась")
	}

	createUsers := `
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            email TEXT NOT NULL UNIQUE,
            username TEXT NOT NULL,
            password TEXT NOT NULL
        );
    `
	_, err = db.Exec(createUsers)
	if err != nil {
		panic(err)
	}

	createSessions := `
		CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			user_id INTEGER NOT NULL,
			expires_at DATETIME NOT NULL
		);
	`

	_, err = db.Exec(createSessions)
	if err != nil {
		panic(err)
	}

	createPosts := `
	CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    category TEXT,
    created_at DATETIME NOT NULL,
    count_likes INTEGER NOT NULL DEFAULT 0
	);
	`

	_, err = db.Exec(createPosts)
	if err != nil {
		panic(err)
	}

	createComments := `
    CREATE TABLE IF NOT EXISTS comments (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        post_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        content TEXT NOT NULL,
        created_at DATETIME NOT NULL,
        likes INTEGER NOT NULL DEFAULT 0,
        dislikes INTEGER NOT NULL DEFAULT 0
    );
`
	_, err = db.Exec(createComments)
	if err != nil {
		panic(err)
	}

	return db

}
