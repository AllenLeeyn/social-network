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
	Status       string     `json:"status"`
	CreatedBy    int        `json:"created_by"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedBy    int        `json:"updated_by"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

type FollowingResponse struct {
	LeaderUUID   string     `json:"leader_uuid"`
	LeaderName   string     `json:"leader_name"`
	FollowerUUID string     `json:"follower_uuid"`
	FollowerName string     `json:"follower_name"`
	GroupUUID    string     `json:"group_uuid"`
	Status       string     `json:"status"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

func SelectIDsFromUUIDs(f *Following) error {
	if f.LeaderID == 0 {
		leaderID, err := userModel.SelectUserIDByUUID(f.LeaderUUID)
		if err != nil {
			return fmt.Errorf("user not found")
		}
		f.LeaderID = leaderID
	}

	if f.FollowerID == 0 {
		followerID, err := userModel.SelectUserIDByUUID(f.FollowerUUID)
		if err != nil {
			return fmt.Errorf("user not found")
		}
		f.FollowerID = followerID
	}
	return nil
}

func InsertFollowing(f *Following) error {
	if err := SelectIDsFromUUIDs(f); err != nil {
		return err
	}

	qry := `INSERT INTO following (
				leader_id, follower_id, group_id, status, created_by
			) VALUES (?, ?, 
				(SELECT id FROM groups WHERE uuid = ?), ?, ?);`
	_, err := sqlDB.Exec(qry,
		&f.LeaderID,
		&f.FollowerID,
		&f.GroupUUID,
		&f.Status,
		&f.CreatedBy)
	if err != nil {
		return err
	}
	return nil
}

func SelectStatus(f *Following) (string, error) {
	qry := `SELECT status
			FROM following
			WHERE follower_id = ? AND leader_id = ? AND group_id = ?`

	var status string
	err := sqlDB.QueryRow(qry, f.FollowerID, f.LeaderID, f.GroupID).Scan(&status)
	if err != nil {
		return "", err
	}
	return status, nil
}

func UpdateFollowing(f *Following) error {
	qry := `UPDATE following 
			SET status = ?, updated_by = ?, updated_at = CURRENT_TIMESTAMP
			WHERE leader_id = ? AND follower_id = ? AND group_id = ?;`
	_, err := sqlDB.Exec(qry,
		&f.Status,
		&f.UpdatedBy,
		&f.LeaderID,
		&f.FollowerID,
		&f.GroupID)
	if err != nil {
		return err
	}
	return nil
}

func SelectFollowings(userID int, fStatus string) (*[]FollowingResponse, error) {
	if fStatus != "accepted" {
		fStatus = "requested"
	}
	qry := `SELECT 
				follower.uuid AS follower_uuid, follower.nick_name AS follower_name,
				leader.uuid AS leader_uuid, leader.nick_name AS leader_name,
				f.status, f.created_at
			FROM following f
				JOIN users follower ON f.follower_id = follower.id
				JOIN users leader   ON f.leader_id   = leader.id
			WHERE (f.leader_id = ? OR f.follower_id = ?)
				AND f.status = ?
				AND f.group_id = 0
			ORDER BY f.created_at DESC;`

	rows, err := sqlDB.Query(qry, userID, userID, fStatus)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []FollowingResponse
	for rows.Next() {
		var fr FollowingResponse
		err := rows.Scan(
			&fr.FollowerUUID, &fr.FollowerName,
			&fr.LeaderUUID, &fr.LeaderName,
			&fr.Status, &fr.CreatedAt)
		if err != nil {
			return nil, err
		}
		requests = append(requests, fr)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &requests, nil
}
