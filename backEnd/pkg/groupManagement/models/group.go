package models

import (
	"social-network/pkg/utils"
	"time"
)

type Group struct {
	ID           int        `json:"id"`
	UUID         string     `json:"uuid"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Banner_image string     `json:"banner_image"`
	Status       string     `json:"status"`
	CreatedBy    int        `json:"created_by"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedBy    int        `json:"updated_by"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

func InsertGroup(group *Group) (int, string, error) {
	uuid, err := utils.GenerateUuid()
	if err != nil {
		return -1, "", err
	}
	group.UUID = uuid

	qry := `INSERT INTO groups (
				uuid, 
				title, description, banner_image, 
				created_by
			) VALUES (?, ?, ?, ?, ?);`
	result, err := sqlDB.Exec(qry,
		group.UUID,
		group.Title, group.Description, group.Banner_image,
		group.CreatedBy)
	if err != nil {
		return -1, "", err
	}

	groupId, err := result.LastInsertId()
	if err != nil {
		return -1, "", err
	}
	return int(groupId), group.UUID, err
}

func UpdateUser(group *Group) error {
	updateQuery := `
		UPDATE users
		SET title = ?,	description = ?, banner_image =?,
			Status = ?, updated_by = ?, update_at = CURRENT_TIMESTAMP,
		WHERE uuid = ?;`

	_, err := sqlDB.Exec(updateQuery,
		group.Title, group.Description, group.Banner_image,
		group.Status, group.UpdatedBy,
		group.UUID,
	)

	return err
}
