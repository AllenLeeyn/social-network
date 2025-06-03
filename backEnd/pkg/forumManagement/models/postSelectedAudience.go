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

func InsertPostSelectedAudience(postId int, selectedAudienceUserUUIDS []string, userId int, tx *sql.Tx) error {
	if len(selectedAudienceUserUUIDS) > 0 {
		query := `INSERT INTO post_selected_audience (post_id, user_id, created_by) VALUES `
		values := make([]any, 0, len(selectedAudienceUserUUIDS)*3)

		for i, uuid := range selectedAudienceUserUUIDS {
			if i > 0 {
				query += ", "
			}
			query += `(?, (SELECT id FROM users WHERE uuid = ?), ?)`
			values = append(values, postId, uuid, userId)
		}
		query += `ON CONFLICT(post_id, user_id) DO UPDATE SET
					status = 'enable',
					updated_at = CURRENT_TIMESTAMP,
					updated_by = excluded.created_by;`

		// Execute the bulk insert query
		_, err := tx.Exec(query, values...)
		if err != nil {
			tx.Rollback() // Rollback on error
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
