package models

import (
	"fmt"
	followingModel "social-network/pkg/followingManagement/models"
)

type memberView followingModel.FollowingView

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

func SelectGroupMembers(tgtUUID, mStatus string) (*[]memberView, error) {
	if mStatus != "accepted" && mStatus != "requested" {
		return nil, fmt.Errorf("invalid status")
	}

	mStatus2 := ""
	if mStatus == "requested" {
		mStatus2 = "invited"
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

	rows, err := sqlDB.Query(qry, mStatus, mStatus2, tgtUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mViews []memberView
	for rows.Next() {
		var mv memberView
		err := rows.Scan(
			&mv.FollowerUUID, &mv.FollowerName,
			&mv.Status, &mv.CreatedAt)
		if err != nil {
			return nil, err
		}
		mViews = append(mViews, mv)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &mViews, nil
}
