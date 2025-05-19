package models

import (
	followingModel "social-network/pkg/followingManagement/models"
)

type Following followingModel.Following

var InsertGroupMember = followingModel.InsertFollowing

func UpdateGroupMember() {}

func SelectGroupMembers() {}

func IsGroupMember(groupUUID string, userID int) bool {
	qry := `SELECT g.id
			FROM groups g
			LEFT JOIN following f ON g.id = f.group_id
			WHERE g.uuid = ? AND f.follower_id = ?`

	var groupID int
	err := sqlDB.QueryRow(qry, groupUUID, userID).Scan(&groupID)
	return err == nil
}
