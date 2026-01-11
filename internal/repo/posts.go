package repo

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"forum/internal/models"
)

type PostCardsFilter struct {
	UserID       int
	CategoryID   int
	MineOnly     bool
	LikedOnly    bool
	CommentLimit int
}

func GetPostCards(db *sql.DB, filter PostCardsFilter) ([]models.PostCard, error) {
	commentLimit := filter.CommentLimit
	if commentLimit <= 0 {
		commentLimit = 3
	}

	query := `
    SELECT
        p.id, p.title, p.content,
        cat.names,
        u.username,
        COALESCE((SELECT SUM(CASE WHEN value = 1 THEN 1 ELSE 0 END) FROM post_reactions pr WHERE pr.post_id = p.id), 0),
        COALESCE((SELECT SUM(CASE WHEN value = -1 THEN 1 ELSE 0 END) FROM post_reactions pr WHERE pr.post_id = p.id), 0),
        cm.id,
        cu.username,
        cm.content
    FROM posts p
    JOIN (
        SELECT pc.post_id, GROUP_CONCAT(c.name, ', ') AS names
        FROM post_categories pc
        JOIN categories c ON c.id = pc.category_id
        GROUP BY pc.post_id
    ) cat ON cat.post_id = p.id
    JOIN users u ON u.id = p.user_id
    LEFT JOIN (
        SELECT c.*, ROW_NUMBER() OVER (PARTITION BY c.post_id ORDER BY c.created_at DESC) AS rn
        FROM comments c
    ) cm ON cm.post_id = p.id AND cm.rn <= ?
    LEFT JOIN users cu ON cu.id = cm.user_id
    `

	var joins []string
	var conditions []string
	args := []any{commentLimit}

	if filter.LikedOnly {
		joins = append(joins, "JOIN post_reactions r ON r.post_id = p.id AND r.user_id = ? AND r.value = 1")
		args = append(args, filter.UserID)
	}
	if filter.MineOnly {
		conditions = append(conditions, "p.user_id = ?")
		args = append(args, filter.UserID)
	}
	if filter.CategoryID != 0 {
		joins = append(joins, "JOIN post_categories pcfilter ON pcfilter.post_id = p.id AND pcfilter.category_id = ?")
		args = append(args, filter.CategoryID)
	}

	if len(joins) > 0 {
		query += "\n" + strings.Join(joins, "\n")
	}
	if len(conditions) > 0 {
		query += "\nWHERE " + strings.Join(conditions, " AND ")
	}

	query += "\nORDER BY p.created_at DESC, cm.created_at DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cards := make([]models.PostCard, 0)
	index := make(map[int]*models.PostCard)

	for rows.Next() {
		var (
			postID          int
			title           string
			content         string
			categoryName    string
			authorName      string
			likes           int
			dislikes        int
			commentID       sql.NullInt64
			commentAuthor   sql.NullString
			commentContent  sql.NullString
		)

		if err := rows.Scan(
			&postID,
			&title,
			&content,
			&categoryName,
			&authorName,
			&likes,
			&dislikes,
			&commentID,
			&commentAuthor,
			&commentContent,
		); err != nil {
			return nil, err
		}

		card, ok := index[postID]
		if !ok {
			newCard := models.PostCard{
				ID:           postID,
				Title:        title,
				Content:      content,
				CategoryName: categoryName,
				AuthorName:   authorName,
				Likes:        likes,
				Dislikes:     dislikes,
			}
			cards = append(cards, newCard)
			card = &cards[len(cards)-1]
			index[postID] = card
		}

		if commentID.Valid {
			card.Comments = append(card.Comments, models.CommentCard{
				AuthorName: commentAuthor.String,
				Content:    commentContent.String,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}

func CreatePost(db *sql.DB, userID int, title string, content string, categoryIDs []int) (int, error) {
	if len(categoryIDs) == 0 {
		return 0, errors.New("category list is empty")
	}
	query := `INSERT INTO posts (user_id, title, content, category_id, created_at) VALUES (?, ?, ?, ?, ?)`

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	res, err := tx.Exec(query, userID, title, content, categoryIDs[0], time.Now())
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	postID, err := res.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	for _, categoryID := range categoryIDs {
		_, err = tx.Exec(`INSERT OR IGNORE INTO post_categories (post_id, category_id) VALUES (?, ?)`, postID, categoryID)
		if err != nil {
			_ = tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return int(postID), nil
}

func GetAllPosts(db *sql.DB) ([]models.PostView, error) {
	query := `
    SELECT 
        p.id, p.user_id, p.title, p.content,
        p.category_id, c.name,
        p.created_at, p.count_likes
    FROM posts p
    JOIN categories c ON c.id = p.category_id
    ORDER BY p.created_at DESC
    `
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.PostView
	for rows.Next() {
		var p models.PostView
		if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CategoryID, &p.CategoryName, &p.CreatedAt, &p.CountLikes); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func GetPostsByCategory(db *sql.DB, categoryID int) ([]models.PostView, error) {
	query := `
    SELECT 
        p.id, p.user_id, p.title, p.content,
        p.category_id, c.name,
        p.created_at, p.count_likes
    FROM posts p
    JOIN categories c ON c.id = p.category_id
    WHERE p.category_id = ?
    ORDER BY p.created_at DESC
    `
	rows, err := db.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.PostView
	for rows.Next() {
		var p models.PostView
		if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CategoryID, &p.CategoryName, &p.CreatedAt, &p.CountLikes); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func GetPostByID(db *sql.DB, postID int) (*models.PostView, error) {
	query := `
    SELECT 
        p.id, p.user_id, p.title, p.content,
        p.category_id, c.name,
        p.created_at, p.count_likes
    FROM posts p
    JOIN categories c ON c.id = p.category_id
    WHERE p.id = ?
    LIMIT 1
    `
	row := db.QueryRow(query, postID)

	var p models.PostView
	err := row.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CategoryID, &p.CategoryName, &p.CreatedAt, &p.CountLikes)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func GetPostCardWithComments(db *sql.DB, postID int) (*models.PostCardWithComments, error) {
	query := `
    SELECT
        p.id, p.title, p.content,
        cat.names,
        u.username,
        COALESCE((SELECT SUM(CASE WHEN value = 1 THEN 1 ELSE 0 END) FROM post_reactions pr WHERE pr.post_id = p.id), 0),
        COALESCE((SELECT SUM(CASE WHEN value = -1 THEN 1 ELSE 0 END) FROM post_reactions pr WHERE pr.post_id = p.id), 0),
        cm.id,
        cm.user_id,
        cu.username,
        cm.content,
        cm.created_at,
        COALESCE((SELECT SUM(CASE WHEN value = 1 THEN 1 ELSE 0 END) FROM comment_reactions cr WHERE cr.comment_id = cm.id), 0),
        COALESCE((SELECT SUM(CASE WHEN value = -1 THEN 1 ELSE 0 END) FROM comment_reactions cr WHERE cr.comment_id = cm.id), 0)
    FROM posts p
    JOIN (
        SELECT pc.post_id, GROUP_CONCAT(c.name, ', ') AS names
        FROM post_categories pc
        JOIN categories c ON c.id = pc.category_id
        GROUP BY pc.post_id
    ) cat ON cat.post_id = p.id
    JOIN users u ON u.id = p.user_id
    LEFT JOIN comments cm ON cm.post_id = p.id
    LEFT JOIN users cu ON cu.id = cm.user_id
    WHERE p.id = ?
    ORDER BY cm.created_at DESC
    `

	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var post *models.PostCardWithComments

	for rows.Next() {
		var (
			id              int
			title           string
			content         string
			categoryName    string
			authorName      string
			likes           int
			dislikes        int
			commentID       sql.NullInt64
			commentUserID   sql.NullInt64
			commentAuthor   sql.NullString
			commentContent  sql.NullString
			commentCreated  sql.NullTime
			commentLikes    sql.NullInt64
			commentDislikes sql.NullInt64
		)

		if err := rows.Scan(
			&id,
			&title,
			&content,
			&categoryName,
			&authorName,
			&likes,
			&dislikes,
			&commentID,
			&commentUserID,
			&commentAuthor,
			&commentContent,
			&commentCreated,
			&commentLikes,
			&commentDislikes,
		); err != nil {
			return nil, err
		}

		if post == nil {
			post = &models.PostCardWithComments{
				ID:           id,
				Title:        title,
				Content:      content,
				CategoryName: categoryName,
				AuthorName:   authorName,
				Likes:        likes,
				Dislikes:     dislikes,
			}
		}

		if commentID.Valid {
			comment := models.CommentView{
				ID:         int(commentID.Int64),
				UserID:     int(commentUserID.Int64),
				AuthorName: commentAuthor.String,
				Content:    commentContent.String,
				CreatedAt:  commentCreated.Time,
				Likes:      int(commentLikes.Int64),
				Dislikes:   int(commentDislikes.Int64),
			}
			post.Comments = append(post.Comments, comment)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if post == nil {
		return nil, sql.ErrNoRows
	}

	return post, nil
}

func GetMyPosts(db *sql.DB, userID int) ([]models.PostView, error) {
	query := `
    SELECT 
        p.id, p.user_id, p.title, p.content,
        p.category_id, c.name,
        p.created_at, p.count_likes
    FROM posts p
    JOIN categories c ON c.id = p.category_id
    WHERE p.user_id = ?
    ORDER BY p.created_at DESC
    `
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.PostView
	for rows.Next() {
		var p models.PostView
		if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CategoryID, &p.CategoryName, &p.CreatedAt, &p.CountLikes); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func GetLikedPosts(db *sql.DB, userID int) ([]models.PostView, error) {
	query := `
    SELECT 
        p.id, p.user_id, p.title, p.content,
        p.category_id, c.name,
        p.created_at, p.count_likes
    FROM posts p
    JOIN categories c ON c.id = p.category_id
    JOIN post_reactions r ON r.post_id = p.id
    WHERE r.user_id = ? AND r.value = 1
    ORDER BY p.created_at DESC
    `
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.PostView
	for rows.Next() {
		var p models.PostView
		if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CategoryID, &p.CategoryName, &p.CreatedAt, &p.CountLikes); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func PostExists(db *sql.DB, postID int) (bool, error) {
	row := db.QueryRow(`SELECT 1 FROM posts WHERE id = ? LIMIT 1`, postID)
	var one int
	err := row.Scan(&one)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
