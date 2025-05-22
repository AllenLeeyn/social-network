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
	BannerImage  string     `json:"banner_image"`
	MembersCount int        `json:"members_count"`
	Status       string     `json:"status"`
	CreatedBy    int        `json:"created_by"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedBy    int        `json:"updated_by"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

type groupView struct {
	UUID         string `json:"uuid"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	BannerImage  string `json:"banner_image"`
	MembersCount int    `json:"members_count"`
	CreatorName  string `json:"creator_name"`
	CreatorUUID  string `json:"creator_uuid"`
	Status       string `json:"status"`
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
		group.Title, group.Description, group.BannerImage,
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

func UpdateGroup(group *Group) error {
	updateQuery := `
		UPDATE groups
		SET title = ?,	description = ?, banner_image =?,
			updated_by = ?, updated_at = CURRENT_TIMESTAMP
		WHERE uuid = ?;`

	_, err := sqlDB.Exec(updateQuery,
		group.Title, group.Description, group.BannerImage,
		group.UpdatedBy,
		group.UUID,
	)
	return err
}

func SelectGroups(userUUID string, joinedOnly bool) (*[]groupView, error) {
	joinedOnlyQry := ``
	if joinedOnly {
		joinedOnlyQry = ` AND f.status = 'accepted'`
	}
	qry := `SELECT
				g.uuid AS group_uuid, g.title,
				g.description, g.banner_image, g.members_count,
				u.nick_name AS creator_name, u.uuid AS creator_uuid,
				COALESCE(f.status, '') as status
			FROM groups g
			JOIN users u ON g.created_by = u.id
			LEFT JOIN users u2 ON u2.uuid = ?
			LEFT JOIN following f ON f.group_id = g.id
				AND f.follower_id = u2.id
			WHERE g.status = 'enable' AND g.id != 0` + joinedOnlyQry + `
			ORDER BY g.created_at DESC;`

	rows, err := sqlDB.Query(qry, userUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []groupView
	for rows.Next() {
		var g groupView
		err := rows.Scan(
			&g.UUID, &g.Title,
			&g.Description, &g.BannerImage, &g.MembersCount,
			&g.CreatorName, &g.CreatorUUID,
			&g.Status)
		if err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &groups, nil
}

func SelectGroup(userID int, groupUUID string) (*groupView, error) {
	qry := `SELECT
				g.uuid AS group_uuid, g.title,
				g.description, g.banner_image, g.members_count,
				u.nick_name AS creator_name, u.uuid AS creator_uuid,
				COALESCE(f.status, '') as status
			FROM groups g
			JOIN users u ON g.created_by = u.id
			LEFT JOIN following f 
				ON f.group_id = g.id 
				AND f.follower_id = ? 
			WHERE g.status = 'enable' AND g.id != 0 AND g.uuid = ?;`
	var g groupView
	err := sqlDB.QueryRow(qry, userID, groupUUID).Scan(
		&g.UUID, &g.Title,
		&g.Description, &g.BannerImage, &g.MembersCount,
		&g.CreatorName, &g.CreatorUUID,
		&g.Status)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func SelectGroupIDcreatedByfromUUID(groupUUID string) (int, int, error) {
	var groupID, createdBy int
	err := sqlDB.QueryRow(`SELECT id, created_by FROM groups WHERE uuid = ?`, groupUUID).
		Scan(&groupID, &createdBy)
	return groupID, createdBy, err
}
