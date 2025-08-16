package models

import (
	"database/sql"
	"time"
)

// Post represents a forum post.
type Post struct {
	ID         int
	UserID     int
	Username   string
	Title      string
	Content    string
	Categories []string
	CreatedAt  time.Time
	Likes      int
	Dislikes   int
}

// CreatePost creates a new post with the given categories.
func CreatePost(db *sql.DB, userID int, title, content string, categories []string) (int, error) {
	res, err := db.Exec("INSERT INTO posts(user_id, title, content) VALUES (?,?,?)", userID, title, content)
	if err != nil {
		return 0, err
	}
	postID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	for _, cat := range categories {
		if cat == "" {
			continue
		}
		_, err = db.Exec("INSERT OR IGNORE INTO categories(name) VALUES (?)", cat)
		if err != nil {
			return 0, err
		}
		var catID int
		if err = db.QueryRow("SELECT id FROM categories WHERE name=?", cat).Scan(&catID); err != nil {
			return 0, err
		}
		_, err = db.Exec("INSERT INTO post_categories(post_id, category_id) VALUES (?,?)", postID, catID)
		if err != nil {
			return 0, err
		}
	}
	return int(postID), nil
}

// GetPost retrieves a single post by ID with categories and like counts.
func GetPost(db *sql.DB, id int) (*Post, error) {
	p := &Post{}
	row := db.QueryRow(`SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at,
        IFNULL(SUM(CASE WHEN pl.value=1 THEN 1 END),0) as likes,
        IFNULL(SUM(CASE WHEN pl.value=-1 THEN 1 END),0) as dislikes
        FROM posts p
        JOIN users u ON p.user_id=u.id
        LEFT JOIN post_likes pl ON p.id=pl.post_id
        WHERE p.id=? GROUP BY p.id`, id)
	if err := row.Scan(&p.ID, &p.UserID, &p.Username, &p.Title, &p.Content, &p.CreatedAt, &p.Likes, &p.Dislikes); err != nil {
		return nil, err
	}
	p.Categories, _ = getCategories(db, p.ID)
	return p, nil
}

// GetPosts returns posts optionally filtered by category, user or liked by user.
func GetPosts(db *sql.DB, category string, userID, likedBy int) ([]Post, error) {
	var rows *sql.Rows
	var err error
	baseQuery := `SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at,
        IFNULL(SUM(CASE WHEN pl.value=1 THEN 1 END),0) as likes,
        IFNULL(SUM(CASE WHEN pl.value=-1 THEN 1 END),0) as dislikes
        FROM posts p
        JOIN users u ON p.user_id=u.id
        LEFT JOIN post_likes pl ON p.id=pl.post_id`
	where := ""
	args := []interface{}{}
	if category != "" {
		baseQuery += " JOIN post_categories pc ON p.id=pc.post_id JOIN categories c ON pc.category_id=c.id"
		where = " WHERE c.name = ?"
		args = append(args, category)
	}
	if userID != 0 {
		if where == "" {
			where = " WHERE p.user_id = ?"
		} else {
			where += " AND p.user_id = ?"
		}
		args = append(args, userID)
	}
	if likedBy != 0 {
		if where == "" {
			where = " WHERE pl.user_id = ? AND pl.value = 1"
		} else {
			where += " AND pl.user_id = ? AND pl.value = 1"
		}
		args = append(args, likedBy)
	}
	query := baseQuery + where + " GROUP BY p.id ORDER BY p.created_at DESC"
	rows, err = db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		p := Post{}
		if err := rows.Scan(&p.ID, &p.UserID, &p.Username, &p.Title, &p.Content, &p.CreatedAt, &p.Likes, &p.Dislikes); err != nil {
			return nil, err
		}
		p.Categories, _ = getCategories(db, p.ID)
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func getCategories(db *sql.DB, postID int) ([]string, error) {
	rows, err := db.Query(`SELECT c.name FROM categories c JOIN post_categories pc ON c.id=pc.category_id WHERE pc.post_id=?`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cats []string
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}

// AddComment is in comment.go etc
