package models

import (
	followingModel "social-network/pkg/followingManagement/models"
)

type Following followingModel.Following

var InsertGroupMember = followingModel.InsertFollowing

func UpdateGroupMember() {}

func SelectGroupMembers() {}

func IsGroupMember(groupUUID string, userID int) bool {
	qry := `SELECT 1
			FROM groups g
			LEFT JOIN following f ON g.id = f.group_id
			WHERE g.uuid = ? AND f.follower_id = ?
			LIMIT 1;`

	var exists int
	err := sqlDB.QueryRow(qry, groupUUID, userID).Scan(&exists)
	return err == nil
}
