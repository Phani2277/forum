package repo

import (
	"database/sql"

	"forum/internal/models"
)

func SeedCategories(db *sql.DB) error {
	_, err := db.Exec(`
		INSERT OR IGNORE INTO categories (name) VALUES
		('Математика'),
		('Рофл'),
		('Новости'),
		('Игры');
	`)
	return err
}

func GetAllCategories(db *sql.DB) ([]models.Category, error) {
	rows, err := db.Query(`SELECT id, name FROM categories ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func CategoryExists(db *sql.DB, categoryID int) (bool, error) {
	row := db.QueryRow(`SELECT 1 FROM categories WHERE id = ? LIMIT 1`, categoryID)
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
