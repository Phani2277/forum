package models

import "database/sql"

// SetPostLike sets the like value (1 or -1) for a post by a user. value=0 removes like.
func SetPostLike(db *sql.DB, userID, postID, value int) error {
	var existing int
	err := db.QueryRow("SELECT value FROM post_likes WHERE user_id=? AND post_id=?", userID, postID).Scan(&existing)
	if err == sql.ErrNoRows {
		if value == 0 {
			return nil
		}
		_, err = db.Exec("INSERT INTO post_likes(user_id, post_id, value) VALUES (?,?,?)", userID, postID, value)
		return err
	}
	if err != nil {
		return err
	}
	if value == 0 {
		_, err = db.Exec("DELETE FROM post_likes WHERE user_id=? AND post_id=?", userID, postID)
	} else {
		_, err = db.Exec("UPDATE post_likes SET value=? WHERE user_id=? AND post_id=?", value, userID, postID)
	}
	return err
}

// SetCommentLike sets the like value for a comment by a user.
func SetCommentLike(db *sql.DB, userID, commentID, value int) error {
	var existing int
	err := db.QueryRow("SELECT value FROM comment_likes WHERE user_id=? AND comment_id=?", userID, commentID).Scan(&existing)
	if err == sql.ErrNoRows {
		if value == 0 {
			return nil
		}
		_, err = db.Exec("INSERT INTO comment_likes(user_id, comment_id, value) VALUES (?,?,?)", userID, commentID, value)
		return err
	}
	if err != nil {
		return err
	}
	if value == 0 {
		_, err = db.Exec("DELETE FROM comment_likes WHERE user_id=? AND comment_id=?", userID, commentID)
	} else {
		_, err = db.Exec("UPDATE comment_likes SET value=? WHERE user_id=? AND comment_id=?", value, userID, commentID)
	}
	return err
}
