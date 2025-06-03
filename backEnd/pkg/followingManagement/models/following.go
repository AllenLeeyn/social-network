package models

import (
	"fmt"
	"time"

	userModel "social-network/pkg/userManagement/models"
)

type Following struct {
	LeaderUUID   string `json:"leader_uuid"`
	LeaderID     int
	FollowerUUID string `json:"follower_uuid"`
	FollowerID   int
	GroupUUID    string `json:"group_uuid"`
	GroupID      int
	Type         string
	Status       string     `json:"status"`
	CreatedBy    int        `json:"created_by"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedBy    int        `json:"updated_by"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

type FollowingView struct {
	LeaderUUID   string     `json:"leader_uuid,omitempty"`
	LeaderName   string     `json:"leader_name,omitempty"`
	FollowerUUID string     `json:"follower_uuid"`
	FollowerName string     `json:"follower_name"`
	GroupUUID    string     `json:"group_uuid"`
	Status       string     `json:"status"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

var publicGroupUUID = "00000000-0000-0000-0000-000000000000"

func IsFollower(userID int, tgtUUID string) bool {
	qry := `SELECT 1
			FROM following f
			JOIN users u ON u.id = f.leader_id
			WHERE u.uuid = ? 
				AND f.follower_id = ? 
				AND f.group_id = 0 
				AND f.status = 'accepted'
			LIMIT 1;`

	var exists int
	err := sqlDB.QueryRow(qry, tgtUUID, userID).Scan(&exists)
	return err == nil
}

func IsLeader(userID int, tgtUUID string) bool {
	qry := `SELECT 1
			FROM following f
			JOIN users u ON u.id = f.follower_id
			WHERE u.uuid = ? 
				AND f.leader_id = ? 
				AND f.group_id = 0 
				AND f.status = 'accepted'
			LIMIT 1;`

	var exists int
	err := sqlDB.QueryRow(qry, tgtUUID, userID).Scan(&exists)
	return err == nil
}

func SelectIDsFromUUIDs(f *Following) error {
	if f.LeaderID == 0 && f.GroupUUID == "" {
		leaderID, err := userModel.SelectUserIDByUUID(f.LeaderUUID)
		if err != nil {
			return fmt.Errorf("leader not found")
		}
		f.LeaderID = leaderID
	}

	if f.FollowerID == 0 {
		followerID, err := userModel.SelectUserIDByUUID(f.FollowerUUID)
		if err != nil {
			return fmt.Errorf("follower not found")
		}
		f.FollowerID = followerID
	}

	return nil
}

func SelectStatus(f *Following) (string, error) {
	qry := `SELECT status
			FROM following
			WHERE follower_id = ? AND leader_id = ? AND group_id = ?`

	var status string
	err := sqlDB.QueryRow(qry, f.FollowerID, f.LeaderID, f.GroupID).Scan(&status)
	return status, checkErrNoRows(err)
}

func SelectFollowings(tgtUUID, fStatus string) (*[]FollowingView, error) {
	if fStatus != "accepted" && fStatus != "requested" {
		return nil, fmt.Errorf("invalid status")
	}

	qry := `SELECT 
				follower.uuid, follower.nick_name,
				leader.uuid, leader.nick_name,
				f.status, f.created_at
			FROM following f
				JOIN users follower ON f.follower_id = follower.id
				JOIN users leader   ON f.leader_id = leader.id
			WHERE (follower.uuid = ? OR leader.uuid = ?)
				AND f.status = ?
				AND f.group_id = 0
			ORDER BY f.created_at DESC;`

	rows, err := sqlDB.Query(qry, tgtUUID, tgtUUID, fStatus)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fViews []FollowingView
	for rows.Next() {
		var fv FollowingView
		err := rows.Scan(
			&fv.FollowerUUID, &fv.FollowerName,
			&fv.LeaderUUID, &fv.LeaderName,
			&fv.Status, &fv.CreatedAt)
		if err != nil {
			return nil, err
		}
		fViews = append(fViews, fv)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &fViews, nil
}

func InsertFollowing(f *Following) error {
	if f.GroupID == 0 {
		f.GroupUUID = publicGroupUUID
	}

	qry := `INSERT INTO following (
				leader_id, follower_id, type, group_id, status, created_by
			) VALUES (?, ?, ?,
				(SELECT id FROM groups WHERE uuid = ?), ?, ?);`
	_, err := sqlDB.Exec(qry,
		f.LeaderID, f.FollowerID, f.Type, f.GroupUUID, f.Status, f.CreatedBy)
	if err != nil {
		return err
	}
	return nil
}

func UpdateFollowing(f *Following) error {
	qry := `UPDATE following 
			SET status = ?, updated_by = ?, updated_at = CURRENT_TIMESTAMP
			WHERE leader_id = ? AND follower_id = ? AND group_id = ?;`
	_, err := sqlDB.Exec(qry,
		f.Status, f.UpdatedBy, f.LeaderID, f.FollowerID, f.GroupID)
	if err != nil {
		return err
	}
	return nil
}
