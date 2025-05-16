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
	GroupID      int        `json:"group_id"`
	Status       string     `json:"status"`
	CreatedBy    int        `json:"created_by"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedBy    int        `json:"updated_by"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

func InsertFollowing(following *Following) error {
	if following.LeaderID == 0 {
		leaderID, err := userModel.SelectUserIDByUUID(following.LeaderUUID)
		if err != nil {
			return fmt.Errorf("user not found")
		}
		following.LeaderID = leaderID
	}

	if following.FollowerID == 0 {
		followerID, err := userModel.SelectUserIDByUUID(following.FollowerUUID)
		if err != nil {
			return fmt.Errorf("user not found")
		}
		following.FollowerID = followerID
	}

	qry := `INSERT INTO following (
				leader_id, follower_id, group_id, status, created_by
			) VALUES (?, ?, ?, ?, ?);`
	_, err := sqlDB.Exec(qry,
		&following.LeaderID,
		&following.FollowerID,
		&following.GroupID,
		&following.Status,
		&following.CreatedBy)
	if err != nil {
		return err
	}
	return nil
}
