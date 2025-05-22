package models

import (
	followingModel "social-network/pkg/followingManagement/models"
)

var InsertGroupMember = followingModel.InsertFollowing

func IsGroupMember(groupUUID string, userID int) bool {
	qry := `SELECT 1
			FROM groups g
			LEFT JOIN following f ON g.id = f.group_id
			WHERE g.uuid = ? 
				AND f.follower_id = ?
				AND f.status = 'accepted'
			LIMIT 1;`

	var exists int
	err := sqlDB.QueryRow(qry, groupUUID, userID).Scan(&exists)
	return err == nil
}

func SelectGroupMembers(tgtUUID, fStatus string) (
	*[]followingModel.FollowingResponse, error) {
	fStatus2 := ""
	if fStatus != "accepted" {
		fStatus, fStatus2 = "requested", "invited"
	}
	qry := `SELECT 
				follower.uuid AS follower_uuid, follower.nick_name AS follower_name,
				f.status, f.created_at
			FROM following f
				JOIN users follower ON f.follower_id = follower.id
				JOIN groups g ON g.id = f.group_id
			WHERE f.status IN (?, ?)
				AND g.uuid = ?
			ORDER BY f.created_at DESC;`

	rows, err := sqlDB.Query(qry, fStatus, fStatus2, tgtUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fResponses []followingModel.FollowingResponse
	for rows.Next() {
		var fr followingModel.FollowingResponse
		err := rows.Scan(
			&fr.FollowerUUID, &fr.FollowerName,
			&fr.Status, &fr.CreatedAt)
		if err != nil {
			return nil, err
		}
		fResponses = append(fResponses, fr)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &fResponses, nil
}
