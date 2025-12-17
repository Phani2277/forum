package main

import "database/sql"

type Category struct {
	ID   int
	Name string
}

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

func GetAllCategories(db *sql.DB) ([]Category, error) {
	rows, err := db.Query(`SELECT id, name FROM categories ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}
