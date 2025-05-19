package models

import (
	"errors"
	"fmt"
	"log"
	userManagementModels "social-network/pkg/userManagement/models"
	"time"
)

type PostFeedback struct {
	ID        int                       `json:"id"`
	Rating    int                       `json:"rating"`
	ParentId  int                       `json:"parent_id"`
	UserId    int                       `json:"user_id"`
	Status    string                    `json:"status"`
	CreatedAt time.Time                 `json:"created_at"`
	UpdatedAt *time.Time                `json:"updated_at"`
	UpdatedBy *int                      `json:"updated_by"`
	User      userManagementModels.User `json:"user"`
	Post      Post                      `json:"post"`
}

func InsertPostFeedback(postFeedback *PostFeedback) (int, error) {
	insertQuery := `INSERT INTO post_feedback (rating, parent_id, user_id) VALUES (?, ?, ?);`
	result, insertErr := sqlDB.Exec(insertQuery, postFeedback.Rating, postFeedback.ParentId, postFeedback.UserId)
	if insertErr != nil {
		// Check if the error is a SQLite constraint violation
		var ErrDuplicatePostFeedback = errors.New("duplicate post feedback")
		if sqliteErr, ok := insertErr.(interface{ ErrorCode() int }); ok {
			// if sqliteErr.ErrorCode() == 19 { // SQLite constraint violation error code
			// 	return -1, sql.ErrNoRows // Return custom error to indicate a duplicate
			// }
			if sqliteErr.ErrorCode() == 19 {
				return -1, ErrDuplicatePostFeedback
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

func UpdatePostFeedback(postFeedback *PostFeedback) error {
	updateQuery := `UPDATE post_feedback
		               SET rating = ?,
			           updated_at = CURRENT_TIMESTAMP,
			           updated_by = ?
		               WHERE parent_id = ?
					   AND user_id = ?;`
	_, updateErr := sqlDB.Exec(updateQuery, postFeedback.Rating, postFeedback.UserId, postFeedback.ParentId, postFeedback.UserId)
	if updateErr != nil {
		// Check if the error is a SQLite constraint violation
		var ErrDuplicatePostFeedback = errors.New("duplicate post rating")
		if sqliteErr, ok := updateErr.(interface{ ErrorCode() int }); ok {
			// if sqliteErr.ErrorCode() == 19 { // SQLite constraint violation error code
			// 	return -1, sql.ErrNoRows // Return custom error to indicate a duplicate
			// }
			if sqliteErr.ErrorCode() == 19 {
				return ErrDuplicatePostFeedback
			}
		}

		return updateErr
	}

	return nil
}

func UpdatePostFeedbackStatus(post_feedback_id int, status string, user_id int) error {
	updateQuery := `UPDATE post_feedback
		               SET status = ?,
			           updated_at = CURRENT_TIMESTAMP,
			           updated_by = ?
		               WHERE id = ?;`
	_, updateErr := sqlDB.Exec(updateQuery, status, user_id, post_feedback_id)
	if updateErr != nil {
		return updateErr
	}
	return nil
}

func ReadAllPostsFeedbacks() ([]PostFeedback, error) {
	// Query the records
	rows, selectError := sqlDB.Query(`
        SELECT 
			pf.id as post_feedback_id, pf.rating, pf.status as post_feedback_status, pf.created_at as post_feedback_created_at, pf.updated_at as post_feedback_updated_at, pf.updated_by as post_feedback_updated_by,
			p.id as post_id, p.status as post_status, p.created_at as post_created_at, p.updated_at as post_updated_at, p.updated_by as post_updated_by,
			u.id as user_id, u.first_name as user_first_name, u.last_name as user_last_name, u.nick_name as user_nick_name, u.email as user_email, IFNULL(u.profile_image, '') as profile_image,
			c.id as category_id, c.name as category_name
		FROM post_feedback pf
			INNER JOIN posts p
				ON pf.parent_id = p.id	
				AND p.status != 'delete'
			INNER JOIN users u
				ON pf.user_id = u.id
				AND u.status != 'delete'
			LEFT JOIN post_categories pc
				ON p.id = pc.post_id
				AND pc.status = 'enable'
			LEFT JOIN categories c
				ON pc.category_id = c.id
				AND c.status = 'enable'
		WHERE pf.status != 'delete'
			;
    `)
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	var postFeedbacks []PostFeedback
	// Map to track postFeedbacks by their ID to avoid duplicates
	postFeedbackMap := make(map[int]*PostFeedback)

	for rows.Next() {
		var postFeedback PostFeedback
		var post Post
		var user userManagementModels.User
		var category Category

		// Scan the post_feedback, post, user, and category data
		err := rows.Scan(
			&postFeedback.ID, &postFeedback.Rating, &postFeedback.Status, &postFeedback.CreatedAt, &postFeedback.UpdatedAt, &postFeedback.UpdatedBy,
			&post.ID, &post.Status, &post.CreatedAt, &post.UpdatedAt, &post.UpdatedBy,
			&user.ID, &user.FirstName, &user.LastName, &user.NickName, &user.Email, &user.ProfileImage,
			&category.ID, &category.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Check if the post_feedback already exists in the postFeedbackMap
		if existingPostFeedback, found := postFeedbackMap[postFeedback.ID]; found {
			// If the post_feedback exists, append the category to the existing post's Categories
			existingPostFeedback.Post.Categories = append(existingPostFeedback.Post.Categories, category)
		} else {
			// If the post_feedback doesn't exist in the map, add it and initialize the Categories field
			postFeedback.Post = post
			postFeedback.User = user
			postFeedback.Post.Categories = []Category{category}
			postFeedbackMap[postFeedback.ID] = &postFeedback
		}
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	// Convert the map of postFeedbacks into a slice
	for _, postFeedback := range postFeedbackMap {
		postFeedbacks = append(postFeedbacks, *postFeedback)
	}

	return postFeedbacks, nil
}

func ReadPostsFeedbacksByUserId(userId int) ([]PostFeedback, error) {
	// Query the records
	rows, selectError := sqlDB.Query(`
        SELECT 
			pf.id as post_feedback_id, pf.rating, pf.status as post_feedback_status, pf.created_at as post_feedback_created_at, pf.updated_at as post_feedback_updated_at, pf.updated_by as post_feedback_updated_by,
			p.id as post_id, p.status as post_status, p.created_at as post_created_at, p.updated_at as post_updated_at, p.updated_by as post_updated_by,
			u.id as user_id, u.first_name as user_first_name, u.last_name as user_last_name, u.nick_name as user_nick_name, u.email as user_email, IFNULL(u.profile_image, '') as profile_image, 
			c.id as category_id, c.name as category_name
		FROM post_feedback pf
			INNER JOIN posts p
				ON pf.parent_id = p.id	
				AND p.status != 'delete'
			INNER JOIN users u
				ON pf.user_id = u.id
				AND u.status != 'delete'
				AND u.id = ?
			LEFT JOIN post_categories pc
				ON p.id = pc.post_id
				AND pc.status = 'enable'
			LEFT JOIN categories c
				ON pc.category_id = c.id
				AND c.status = 'enable'
		WHERE pf.status != 'delete'
			;
    `, userId)
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	var postFeedbacks []PostFeedback
	// Map to track postFeedbacks by their ID to avoid duplicates
	postFeedbackMap := make(map[int]*PostFeedback)

	for rows.Next() {
		var postFeedback PostFeedback
		var post Post
		var user userManagementModels.User
		var category Category

		// Scan the post_feedback, post, user, and category data
		err := rows.Scan(
			&postFeedback.ID, &postFeedback.Rating, &postFeedback.Status, &postFeedback.CreatedAt, &postFeedback.UpdatedAt, &postFeedback.UpdatedBy,
			&post.ID, &post.Status, &post.CreatedAt, &post.UpdatedAt, &post.UpdatedBy,
			&user.ID, &user.FirstName, &user.LastName, &user.NickName, &user.Email, &user.ProfileImage,
			&category.ID, &category.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Check if the post_feedback already exists in the postFeedbackMap
		if existingPostFeedback, found := postFeedbackMap[postFeedback.ID]; found {
			// If the post_feedback exists, append the category to the existing post's Categories
			existingPostFeedback.Post.Categories = append(existingPostFeedback.Post.Categories, category)
		} else {
			// If the post_feedback doesn't exist in the map, add it and initialize the Categories field
			postFeedback.Post = post
			postFeedback.User = user
			postFeedback.Post.Categories = []Category{category}
			postFeedbackMap[postFeedback.ID] = &postFeedback
		}
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	// Convert the map of postFeedbacks into a slice
	for _, postFeedback := range postFeedbackMap {
		postFeedbacks = append(postFeedbacks, *postFeedback)
	}

	return postFeedbacks, nil
}

func ReadPostsFeedbacksByPostId(postId int) ([]PostFeedback, error) {
	// Query the records
	rows, selectError := sqlDB.Query(`
        SELECT 
			pf.id as post_feedback_id, pf.rating, pf.status as post_feedback_status, pf.created_at as post_feedback_created_at, pf.updated_at as post_feedback_updated_at, pf.updated_by as post_feedback_updated_by,
			p.id as post_id, p.status as post_status, p.created_at as post_created_at, p.updated_at as post_updated_at, p.updated_by as post_updated_by,
			u.id as user_id, u.first_name as user_first_name, u.last_name as user_last_name, u.nick_name as user_nick_name, u.email as user_email, IFNULL(u.profile_image, '') as profile_image, 
			c.id as category_id, c.name as category_name
		FROM post_feedback pf
			INNER JOIN posts p
				ON pf.parent_id = p.id	
				AND p.status != 'delete'
				AND p.id = ?
			INNER JOIN users u
				ON pf.user_id = u.id
				AND u.status != 'delete'
			LEFT JOIN post_categories pc
				ON p.id = pc.post_id
				AND pc.status = 'enable'
			LEFT JOIN categories c
				ON pc.category_id = c.id
				AND c.status = 'enable'
		WHERE pf.status != 'delete'
			;
    `, postId)
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	var postFeedbacks []PostFeedback
	// Map to track postFeedbacks by their ID to avoid duplicates
	postFeedbackMap := make(map[int]*PostFeedback)

	for rows.Next() {
		var postFeedback PostFeedback
		var post Post
		var user userManagementModels.User
		var category Category

		// Scan the post_feedback, post, user, and category data
		err := rows.Scan(
			&postFeedback.ID, &postFeedback.Rating, &postFeedback.Status, &postFeedback.CreatedAt, &postFeedback.UpdatedAt, &postFeedback.UpdatedBy,
			&post.ID, &post.Status, &post.CreatedAt, &post.UpdatedAt, &post.UpdatedBy,
			&user.ID, &user.FirstName, &user.LastName, &user.NickName, &user.Email, &user.ProfileImage,
			&category.ID, &category.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Check if the post_feedback already exists in the postFeedbackMap
		if existingPostFeedback, found := postFeedbackMap[postFeedback.ID]; found {
			// If the post_feedback exists, append the category to the existing post's Categories
			existingPostFeedback.Post.Categories = append(existingPostFeedback.Post.Categories, category)
		} else {
			// If the post_feedback doesn't exist in the map, add it and initialize the Categories field
			postFeedback.Post = post
			postFeedback.User = user
			postFeedback.Post.Categories = []Category{category}
			postFeedbackMap[postFeedback.ID] = &postFeedback
		}
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	// Convert the map of postFeedbacks into a slice
	for _, postFeedback := range postFeedbackMap {
		postFeedbacks = append(postFeedbacks, *postFeedback)
	}

	return postFeedbacks, nil
}

func PostHasFeedback(userId int, postID int) int {
	var existingLikeRating int
	feedbackCheckQuery := `SELECT rating
		FROM post_feedback pf
		WHERE pf.user_id = ? AND pf.parent_id = ?
		AND status = 'enable'
	`
	err := sqlDB.QueryRow(feedbackCheckQuery, userId, postID).Scan(&existingLikeRating)

	if err == nil { //it means that post has like or dislike
		return existingLikeRating
	} else {
		return -1000
	}
}
