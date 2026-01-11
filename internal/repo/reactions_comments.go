package repo

import "database/sql"

func GetCommentReactionCounts(db *sql.DB, commentID int) (int, int, error) {
	row := db.QueryRow(`
SELECT
  COALESCE(SUM(CASE WHEN value = 1 THEN 1 ELSE 0 END), 0),
  COALESCE(SUM(CASE WHEN value = -1 THEN 1 ELSE 0 END), 0)
FROM comment_reactions
WHERE comment_id = ?;
`, commentID)

	var likes, dislikes int
	if err := row.Scan(&likes, &dislikes); err != nil {
		return 0, 0, err
	}
	return likes, dislikes, nil
}
