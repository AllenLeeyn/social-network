package models

import (
	"database/sql"
	"time"
)

type PostSelectedAudience struct {
	PostId    sql.NullInt64 `json:"post_id"`
	UserId    sql.NullInt64 `json:"user_id"`
	Status    string        `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	CreatedBy int           `json:"created_by"`
	UpdatedAt time.Time     `json:"updated_at"`
	UpdatedBy int           `json:"updated_by"`
}

func InsertPostSelectedAudience(postId int, selectedAudienceUserIds []int, userId int, tx *sql.Tx) error {
	if len(selectedAudienceUserIds) == 0 {
		return nil
	}

	insertQuery := `INSERT INTO post_selected_audience (post_id, user_id, created_at, created_by) VALUES (?, ?, CURRENT_TIMESTAMP, ?);`
	for _, userId := range selectedAudienceUserIds {
		_, err := tx.Exec(insertQuery, postId, userId, userId)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateStatusPostSelectedAudience(postId int, userId int, status string, tx *sql.Tx) error {
	updateQuery := `UPDATE post_selected_audience
					SET status = ?,
						updated_at = CURRENT_TIMESTAMP,
						updated_by = ?
					WHERE post_id = ?;`
	_, updateErr := tx.Exec(updateQuery, status, userId, postId)
	if updateErr != nil {
		return updateErr
	}

	return nil
}
