package models

import (
	"database/sql"
	"fmt"
	"log"
	"social-network/pkg/dbTools"
)

// Category struct represents the user data model
type Category struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	PostsCount     *int   `json:"posts_count"`
	CommentsCount  *int   `json:"comments_count"`
	PostLikesCount *int   `json:"post_likes_count"`
}

func InsertCategory(db *dbTools.DBContainer, category *Category) (int, error) {
	insertQuery := `INSERT INTO categories (name) VALUES (?);`
	result, insertErr := db.Conn.Exec(insertQuery, category.Name)
	if insertErr != nil {
		// Check if the error is a SQLite constraint violation
		if sqliteErr, ok := insertErr.(interface{ ErrorCode() int }); ok {
			if sqliteErr.ErrorCode() == 19 { // SQLite constraint violation error code
				return -1, sql.ErrNoRows // Return custom error to indicate a duplicate
			}
		}
		return -1, insertErr
	}

	// Retrieve the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return -1, err
	}

	return int(lastInsertID), nil
}

func UpdateCategory(db *dbTools.DBContainer, category *Category, userId int) error {
	updateQuery := `UPDATE categories
					SET name = ?,
						updated_at = CURRENT_TIMESTAMP,
						updated_by = ?
					WHERE id = ?;`
	_, updateErr := db.Conn.Exec(updateQuery, category.Name, userId, category.ID)
	if updateErr != nil {
		// Check if the error is a SQLite constraint violation
		if sqliteErr, ok := updateErr.(interface{ ErrorCode() int }); ok {
			if sqliteErr.ErrorCode() == 19 { // SQLite constraint violation error code
				return sql.ErrNoRows // Return custom error to indicate a duplicate
			}
		}
		return updateErr
	}

	return nil
}

func AdminReadAllCategories(db *dbTools.DBContainer) ([]Category, error) {
	// Query the records
	rows, selectError := db.Conn.Query(`
        SELECT c.id as category_id, c.name as category_name,
			   (SELECT COUNT(DISTINCT p.id) 
			   	FROM post_categories pc
				INNER JOIN posts p
					ON pc.post_id = p.id
				WHERE p.status != 'delete'
				AND pc.status != 'delete'
				AND pc.category_id = c.id
			   ) as posts_count,
			   (SELECT COUNT(DISTINCT com.id) 
			   	FROM post_categories pc
				INNER JOIN posts p
					ON pc.post_id = p.id
				INNER JOIN comments com
					ON com.post_id = p.id
					AND com.status != 'delete'
				WHERE p.status != 'delete'
				AND pc.status != 'delete'
				AND pc.category_id = c.id
			   ) as comments_count,
			   (SELECT COUNT(DISTINCT pl.id) 
			   	FROM post_categories pc
				INNER JOIN posts p
					ON pc.post_id = p.id
				INNER JOIN post_likes pl
					ON pl.post_id = p.id
					AND pl.status != 'delete'
				WHERE p.status != 'delete'
				AND pc.status != 'delete'
				AND pc.category_id = c.id
			   ) as post_likes_count
        FROM categories c
        WHERE c.status != 'delete';
    `)
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	var categories []Category

	for rows.Next() {
		var category Category

		// Scan the data into variables
		err := rows.Scan(
			&category.ID, &category.Name,
			&category.PostsCount, &category.CommentsCount, &category.PostLikesCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Append category to the categories slice
		categories = append(categories, category)
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return categories, nil
}

func ReadAllCategories(db *dbTools.DBContainer) ([]Category, error) {
	// Query the records
	rows, selectError := db.Conn.Query(`
        SELECT c.id as category_id, c.name as category_name
        FROM categories c
        WHERE c.status != 'delete';
    `)
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	var categories []Category

	for rows.Next() {
		var category Category

		// Scan the data into variables
		err := rows.Scan(
			&category.ID, &category.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Append category to the categories slice
		categories = append(categories, category)
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return categories, nil
}

func ReadCategoryById(db *dbTools.DBContainer, categoryId int) (Category, error) {
	// Query the records
	rows, selectError := db.Conn.Query(`
        SELECT c.id as category_id, c.name as category_name
        FROM categories c
        WHERE c.id = ?;
    `, categoryId)
	if selectError != nil {
		return Category{}, selectError
	}
	defer rows.Close()

	// Variable to hold the category and user data
	var category Category

	// Scan the result into variables
	if rows.Next() {
		err := rows.Scan(
			&category.ID, &category.Name,
		)
		if err != nil {
			return Category{}, fmt.Errorf("error scanning row: %v", err)
		}
	} else {
		// If no category found with the given ID
		return Category{}, fmt.Errorf("category with ID %d not found", categoryId)
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return Category{}, fmt.Errorf("row iteration error: %v", err)
	}

	return category, nil
}

func ReadCategoryByName(db *dbTools.DBContainer, categoryName string) (Category, error) {
	// Query the records
	rows, selectError := db.Conn.Query(`
        SELECT c.id as category_id, c.name as category_name, c.color as category_color
        FROM categories c
        WHERE c.name = ?;
    `, categoryName)
	if selectError != nil {
		return Category{}, selectError
	}
	defer rows.Close()

	// Variable to hold the category and user data
	var category Category

	// Scan the result into variables
	if rows.Next() {
		err := rows.Scan(
			&category.ID, &category.Name,
		)
		if err != nil {
			return Category{}, fmt.Errorf("error scanning row: %v", err)
		}
	} else {
		// If no category found with the given Name
		return Category{}, fmt.Errorf("category with Name %v not found", categoryName)
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return Category{}, fmt.Errorf("row iteration error: %v", err)
	}

	return category, nil
}
