package models

import (
	"social-network/pkg/dbTools"
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

func InsertCommentFeedback(db *dbTools.DBContainer, Rating string, parentId int, userId int) error {
	insertQuery := `INSERT INTO comment_feedback (rating, user_id, parent_id) VALUES (?, ?, ?);`
	_, insertErr := db.Conn.Exec(insertQuery, Rating, userId, parentId)
	if insertErr != nil {
		// Check if the error is a SQLite constraint violation
		return insertErr
	}
	return nil
}

func UpdateCommentFeedback(db *dbTools.DBContainer, Rating string, commentFeedback CommentFeedback) error {
	updateQuery := `UPDATE comment_feedback
	SET rating = ?,
		updated_at = CURRENT_TIMESTAMP,
		updated_by = ?
	WHERE id = ?;`
	_, insertErr := db.Conn.Exec(updateQuery, Rating, commentFeedback.UserId, commentFeedback.ID)
	if insertErr != nil {
		// Check if the error is a SQLite constraint violation
		return insertErr
	}
	return nil
}

func UpdateCommentFeedbackStatus(db *dbTools.DBContainer, commentFeedbackId int, status string, user_id int) error {
	updateQuery := `UPDATE comment_feedback
	SET status = ?,
		updated_at = CURRENT_TIMESTAMP,
		updated_by = ?
	WHERE id = ?;`
	_, insertErr := db.Conn.Exec(updateQuery, status, user_id, commentFeedbackId)
	if insertErr != nil {
		// Check if the error is a SQLite constraint violation
		return insertErr
	}
	return nil
}

func ReadAllCommentsFeedbackdByUserId(db *dbTools.DBContainer, userId int, Rating string) ([]CommentFeedback, error) {
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
	rows, insertErr := db.Conn.Query(selectQuery, userId, Rating)
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

func CommentHasFeedback(db *dbTools.DBContainer, userId int, parentIdID int) (int, int) {
	var existingFeedbackId int
	var existingFeedbackRating int
	feedbackCheckQuery := `SELECT id, rating
		FROM comment_feedback cf
		WHERE cf.user_id = ? AND cf.parent_id = ?
		AND status = 'enable'
	`
	err := db.Conn.QueryRow(feedbackCheckQuery, userId, parentIdID).Scan(&existingFeedbackId, &existingFeedbackRating)

	if err == nil { //it means that post has feedback
		return existingFeedbackId, existingFeedbackRating
	} else {
		return -1, 0
	}
}
