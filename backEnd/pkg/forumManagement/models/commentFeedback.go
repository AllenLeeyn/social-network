package models

import (
	"errors"
	"log"
	userManagementModels "social-network/pkg/userManagement/models"
	"time"
)

type CommentFeedback struct {
	ID        int                       `json:"id"`
	Rating    int                       `json:"rating"`
	UserId    int                       `json:"user_id"`
	ParentId  int                       `json:"parent_id"`
	Status    string                    `json:"status"`
	CreatedAt time.Time                 `json:"created_at"`
	UpdatedAt *time.Time                `json:"updated_at"`
	UpdatedBy *int                      `json:"updated_by"`
	Post      Post                      `json:"post"`
	User      userManagementModels.User `json:"user"`
	Comment   Comment                   `json:"comment"`
}

func InsertCommentFeedback(commentFeedback *CommentFeedback) (int, error) {
	insertQuery := `INSERT INTO comment_feedback (rating, parent_id, user_id) VALUES (?, ?, ?);`
	result, insertErr := sqlDB.Exec(insertQuery, commentFeedback.Rating, commentFeedback.ParentId, commentFeedback.UserId)
	if insertErr != nil {
		// Check if the error is a SQLite constraint violation
		var ErrDuplicateCommentFeedback = errors.New("duplicate comment feedback")
		if sqliteErr, ok := insertErr.(interface{ ErrorCode() int }); ok {
			// if sqliteErr.ErrorCode() == 19 { // SQLite constraint violation error code
			// 	return -1, sql.ErrNoRows // Return custom error to indicate a duplicate
			// }
			if sqliteErr.ErrorCode() == 19 {
				return -1, ErrDuplicateCommentFeedback
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

func UpdateCommentFeedback(commentFeedback *CommentFeedback) error {
	updateQuery := `UPDATE comment_feedback
	SET rating = ?,
		updated_at = CURRENT_TIMESTAMP,
		updated_by = ?
	WHERE parent_id = ?
	AND user_id = ?;`
	_, updateErr := sqlDB.Exec(updateQuery, commentFeedback.Rating, commentFeedback.UserId, commentFeedback.ParentId, commentFeedback.UserId)
	if updateErr != nil {
		// Check if the error is a SQLite constraint violation
		var ErrDuplicateCommentFeedback = errors.New("duplicate comment rating")
		if sqliteErr, ok := updateErr.(interface{ ErrorCode() int }); ok {
			// if sqliteErr.ErrorCode() == 19 { // SQLite constraint violation error code
			// 	return -1, sql.ErrNoRows // Return custom error to indicate a duplicate
			// }
			if sqliteErr.ErrorCode() == 19 {
				return ErrDuplicateCommentFeedback
			}
		}

		return updateErr
	}

	return nil
}

func UpdateCommentFeedbackStatus(commentFeedbackId int, status string, user_id int) error {
	updateQuery := `UPDATE comment_feedback
	SET status = ?,
		updated_at = CURRENT_TIMESTAMP,
		updated_by = ?
	WHERE id = ?;`
	_, insertErr := sqlDB.Exec(updateQuery, status, user_id, commentFeedbackId)
	if insertErr != nil {
		// Check if the error is a SQLite constraint violation
		return insertErr
	}
	return nil
}

func ReadAllCommentsFeedbackdByUserId(userId int, Rating string) ([]CommentFeedback, error) {
	selectQuery := `SELECT 
			p.id AS post_id, p.uuid AS post_uuid, p.title AS post_title, p.content AS post_content, 
			p.status AS post_status, p.created_at AS post_created_at, p.updated_at AS post_updated_at, p.updated_by AS post_updated_by,
			u.id AS user_id, u.uuid AS user_uuid, u.nick_name AS user_nick_name, u.first_name as user_first_name, u.last_name as user_last_name, u.type_id AS user_type_id, u.email AS user_email, IFNULL(u.profile_image, '') as profile_image, 
			u.status AS user_status, u.created_at AS user_created_at, u.updated_at AS user_updated_at, u.updated_by AS user_updated_by,

			c.id AS comment_id, c.post_id as comment_post_id, c.content AS comment_content, c.user_id AS comment_user_id, 
			c.status AS comment_status, c.created_at AS comment_created_at, c.updated_at AS comment_updated_at, c.updated_by AS comment_updated_by,


			cf.id AS comment_feedback_id, cf.rating AS comment_feedback_rating, cf.parent_id AS comment_feedback_parent_id, cf.user_id AS comment_feedback_user_id, cf.status AS comment_feedback_status, cf.created_at AS comment_feedback_created_at, cf.updated_at AS comment_feedback_updated_at, cf.updated_by AS comment_feedback_updated_by 
		FROM comment_feedback cf
			INNER JOIN comments c
				ON cf.comment_id = c.id AND cf.user_id = ? AND cf.rating = ? c.status != 'delete' AND cf.status != 'delete' 
			INNER JOIN posts p 
				ON c.post_id = p.id AND p.status != 'delete' 
			INNER JOIN users u 
				ON cf.user_id = u.id AND u.status != 'delete;'		
	`
	rows, insertErr := sqlDB.Query(selectQuery, userId, Rating)
	if insertErr != nil {
		// Check if the error is a SQLite constraint violation
		return nil, insertErr
	}

	var commentFeedbacks []CommentFeedback

	for rows.Next() {
		var commentFeedback CommentFeedback
		var comment Comment
		var user userManagementModels.User
		var post Post
		err := rows.Scan(
			&post.ID,
			&post.UUID,
			&post.Title,
			&post.Content,
			&post.Status,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.UpdatedBy,

			&user.ID,
			&user.UUID,
			&user.NickName,
			&user.FirstName,
			&user.LastName,
			&user.TypeId,
			&user.Email,
			&user.ProfileImage,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.UpdatedBy,

			&comment.ID,
			&comment.PostId,
			&comment.Content,
			&comment.UserId,
			&comment.Status,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.UpdatedBy,

			&commentFeedback.ID,
			&commentFeedback.Rating,
			&commentFeedback.ParentId,
			&commentFeedback.UserId,
			&commentFeedback.Status,
			&commentFeedback.CreatedAt,
			&commentFeedback.UpdatedAt,
			&commentFeedback.UpdatedBy,
		)

		if err != nil {
			return nil, err
		}
		commentFeedback.Comment = comment
		commentFeedback.Post = post
		commentFeedback.User = user

		commentFeedbacks = append(commentFeedbacks, commentFeedback)
	}

	// Check for any errors during the iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return commentFeedbacks, nil

}

func CommentHasFeedback(userId int, parentIdID int) int {
	var existingFeedbackRating int
	feedbackCheckQuery := `SELECT rating
		FROM comment_feedback cf
		WHERE cf.user_id = ? AND cf.parent_id = ?
		AND status = 'enable'
	`
	err := sqlDB.QueryRow(feedbackCheckQuery, userId, parentIdID).Scan(&existingFeedbackRating)

	if err == nil { //it means that post has feedback
		return existingFeedbackRating
	} else {
		return -1000
	}
}
