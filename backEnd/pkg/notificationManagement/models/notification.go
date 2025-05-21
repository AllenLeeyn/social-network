package models

import (
	"database/sql"
	"fmt"
	userManagementModels "social-network/pkg/userManagement/models"
	"time"
)

// Post struct represents the user data model
type Notification struct {
	ID                 int        `json:"id"`
	ToUserId           int        `json:"to_user_id"`
	FromUserId         int        `json:"from_user_id"`
	TargetType         string     `json:"target_type"`
	TargetDetailedType string     `json:"target_detailed_type"`
	TargetId           int        `json:"target_id"`
	Message            string     `json:"message"`
	IsRead             int        `json:"is_read"`
	Data               string     `json:"data"`
	Status             string     `json:"status"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at"`
	UpdatedBy          *int       `json:"updated_by"`

	ToUser   userManagementModels.User `json:"to_user"`   // Embedded to user data
	FromUser userManagementModels.User `json:"from_user"` // Embedded from user data
}

var sqlDB *sql.DB

func Initialize(dbMain *sql.DB) {
	sqlDB = dbMain
}

func InsertNotification(notification *Notification) (int, error) {
	insertQuery := `INSERT INTO notifications (to_user_id, from_user_id, target_type, target_detailed_type, target_id, message, data) VALUES (?, ?, ?, ?, ?, ?, ?);`
	result, insertErr := sqlDB.Exec(insertQuery, notification.ToUser, notification.FromUser, notification.TargetType, notification.TargetDetailedType, notification.TargetId, notification.Message, notification.Data)
	if insertErr != nil {
		return -1, insertErr
	}

	// Retrieve the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(lastInsertID), nil
}

func UpdateNotificationReedStatus(notification_id int, is_read int, user_id int) error {
	updateQuery := `UPDATE notifications
					SET is_read = ?,
						updated_at = CURRENT_TIMESTAMP,
						updated_by = ?
					WHERE id = ?;`
	_, updateErr := sqlDB.Exec(updateQuery, is_read, user_id, notification_id)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

func DeleteNotification(notification_id int, user_id int) error {
	updateQuery := `UPDATE notifications
					SET status = 'delete',
						updated_at = CURRENT_TIMESTAMP,
						updated_by = ?
					WHERE id = ?;`
	_, updateErr := sqlDB.Exec(updateQuery, user_id, notification_id)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

func ReadAllNotifications(to_user_id int) ([]Notification, error) {
	rows, selectError := sqlDB.Query(`
        SELECT n.id as notification_id, n.to_user_id as notification_to_user_id, n.from_user_id as notification_from_user_id, 
			n.target_id as notification_target_id, n.target_type as notification_target_type, n.target_detailed_type as notification_target_detailed_type, 
			n.message as notification_message, n.is_read as notification_is_read, n.data as notification_data,
			n.status as notification_status, n.created_at as notification_created_at, n.updated_at as notification_updated_at, n.updated_by as notification_updated_by,
			to_user.id as to_user_id, to_user.first_name as to_user_first_name, to_user.last_name as to_user_last_name, to_user.nick_name as to_user_nick_name, to_user.email as to_user_email, IFNULL(to_user.profile_image, '') as to_user_profile_image,
			from_user.id as from_user_id, from_user.first_name as from_user_first_name, from_user.last_name as from_user_last_name, from_user.nick_name as from_user_nick_name, from_user.email as from_user_email, IFNULL(from_user.profile_image, '') as from_user_profile_image
		FROM notifications n
			INNER JOIN users to_user
				ON n.to_user_id = to_user.id
				AND to_user.id = ?
			LEFT JOIN users from_user
				ON n.from_user_id = from_user.id
		WHERE n.status != 'delete'
		ORDER BY n.id desc;
    `, to_user_id)
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	var notifications []Notification

	for rows.Next() {
		var notification Notification
		var toUser userManagementModels.User
		var fromUser userManagementModels.User

		// Scan the post and user data
		err := rows.Scan(
			&notification.ID, &notification.ToUserId, &notification.FromUserId,
			&notification.TargetId, &notification.TargetType, &notification.TargetDetailedType,
			&notification.Message, &notification.IsRead, &notification.Data,
			&notification.Status, &notification.CreatedAt, &notification.UpdatedAt, &notification.UpdatedBy,
			&toUser.ID, &toUser.FirstName, &toUser.LastName, &toUser.NickName, &toUser.Email, &toUser.ProfileImage,
			&fromUser.ID, &fromUser.FirstName, &fromUser.LastName, &fromUser.NickName, &fromUser.Email, &fromUser.ProfileImage,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		notification.ToUser = toUser
		notification.FromUser = fromUser
		notifications = append(notifications, notification)
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return notifications, nil
}

func ReadNotificationById(notification_id int, to_user_id int) (Notification, error) {
	rows, selectError := sqlDB.Query(`
        SELECT n.id as notification_id, n.to_user_id as notification_to_user_id, n.from_user_id as notification_from_user_id, 
			n.target_id as notification_target_id, n.target_type as notification_target_type, n.target_detailed_type as notification_target_detailed_type, 
			n.message as notification_message, n.is_read as notification_is_read, n.data as notification_data,
			n.status as notification_status, n.created_at as notification_created_at, n.updated_at as notification_updated_at, n.updated_by as notification_updated_by,
			to_user.id as to_user_id, to_user.first_name as to_user_first_name, to_user.last_name as to_user_last_name, to_user.nick_name as to_user_nick_name, to_user.email as to_user_email, IFNULL(to_user.profile_image, '') as to_user_profile_image,
			from_user.id as from_user_id, from_user.first_name as from_user_first_name, from_user.last_name as from_user_last_name, from_user.nick_name as from_user_nick_name, from_user.email as from_user_email, IFNULL(from_user.profile_image, '') as from_user_profile_image
		FROM notifications n
			INNER JOIN users to_user
				ON n.to_user_id = to_user.id
				AND to_user.id = ?
			LEFT JOIN users from_user
				ON n.from_user_id = from_user.id
		WHERE n.status != 'delete'
			AND n.id = ?
		ORDER BY n.id desc;
    `, to_user_id, notification_id)
	if selectError != nil {
		return Notification{}, selectError
	}
	defer rows.Close()

	var notification Notification

	if rows.Next() {
		var toUser userManagementModels.User
		var fromUser userManagementModels.User

		// Scan the post and user data
		err := rows.Scan(
			&notification.ID, &notification.ToUserId, &notification.FromUserId,
			&notification.TargetId, &notification.TargetType, &notification.TargetDetailedType,
			&notification.Message, &notification.IsRead, &notification.Data,
			&notification.Status, &notification.CreatedAt, &notification.UpdatedAt, &notification.UpdatedBy,
			&toUser.ID, &toUser.FirstName, &toUser.LastName, &toUser.NickName, &toUser.Email, &toUser.ProfileImage,
			&fromUser.ID, &fromUser.FirstName, &fromUser.LastName, &fromUser.NickName, &fromUser.Email, &fromUser.ProfileImage,
		)
		if err != nil {
			return Notification{}, fmt.Errorf("error scanning row: %v", err)
		}

		notification.ToUser = toUser
		notification.FromUser = fromUser
	} else {
		// If no Notification found with the given ID
		return Notification{}, fmt.Errorf("id not found")
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return Notification{}, fmt.Errorf("row iteration error: %v", err)
	}

	return notification, nil
}
