package models

import "database/sql"

// GetCategories returns all category names.
func GetCategories(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT name FROM categories ORDER BY name")
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
