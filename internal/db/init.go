package db

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

	createCategories := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);
	`
	_, err = db.Exec(createCategories)
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
    category_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL,
    count_likes INTEGER NOT NULL DEFAULT 0
	);
	`

	_, err = db.Exec(createPosts)
	if err != nil {
		panic(err)
	}

	createPostCategories := `
	CREATE TABLE IF NOT EXISTS post_categories (
		post_id INTEGER NOT NULL,
		category_id INTEGER NOT NULL,
		PRIMARY KEY (post_id, category_id)
	);
	`
	_, err = db.Exec(createPostCategories)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`INSERT OR IGNORE INTO post_categories (post_id, category_id) SELECT id, category_id FROM posts`)
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

	createCommentReactions := `
	CREATE TABLE IF NOT EXISTS comment_reactions (
		user_id INTEGER NOT NULL,
		comment_id INTEGER NOT NULL,
		value INTEGER NOT NULL,
		PRIMARY KEY (user_id, comment_id)
	);
	`
	_, err = db.Exec(createCommentReactions)
	if err != nil {
		panic(err)
	}

	postReactions := `
	CREATE TABLE IF NOT EXISTS post_reactions (
        user_id INTEGER NOT NULL,
    	post_id INTEGER NOT NULL,
    	value INTEGER NOT NULL,
    	PRIMARY KEY (user_id, post_id)
    )
	`
	_, err = db.Exec(postReactions)
	if err != nil {
		panic(err)
	}

	return db

}
