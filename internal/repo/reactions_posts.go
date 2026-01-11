package repo

import "database/sql"

func GetPostReactionCounts(db *sql.DB, postID int) (int, int, error) {
	row := db.QueryRow(`SELECT
  COALESCE(SUM(CASE WHEN value = 1 THEN 1 ELSE 0 END), 0),
  COALESCE(SUM(CASE WHEN value = -1 THEN 1 ELSE 0 END), 0)
FROM post_reactions
WHERE post_id = ?;`, postID)

	var likes int
	var dislikes int
	if err := row.Scan(&likes, &dislikes); err != nil {
		return 0, 0, err
	}

	return likes, dislikes, nil
}
