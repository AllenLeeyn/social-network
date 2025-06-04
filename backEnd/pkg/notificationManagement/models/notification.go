package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type NotificationFromUser struct {
	ID           int    `json:"id"`
	UUID         string `json:"uuid"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	NickName     string `json:"nick_name"`
	ProfileImage string `json:"profile_image,omitempty"` // Optional field for profile image
}

// Post struct represents the user data model
type Notification struct {
	ID                 int    `json:"id"`
	ToUserId           int    `json:"to_user_id"`
	FromUserId         int    `json:"from_user_id"`
	TargetType         string `json:"target_type"`
	TargetDetailedType string `json:"target_detailed_type"`
	TargetId           int    `json:"target_id"`
	TargetUUID         string
	TargetUUIDForm     string           `json:"target_uuid"`
	Message            string           `json:"message"`
	IsRead             int              `json:"is_read"`
	Data               *json.RawMessage `json:"data"`
	Status             string           `json:"status"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          *time.Time       `json:"updated_at"`
	UpdatedBy          *int             `json:"updated_by"`

	FromUser   NotificationFromUser `json:"from_user"` // Embedded from user data
	ToUserUUID string               `json:"to_user_uuid"`
}

var sqlDB *sql.DB

func Initialize(dbMain *sql.DB) {
	sqlDB = dbMain
}

func InsertNotification(notification *Notification) (int, error) {
	insertQuery := `INSERT INTO notifications (
					to_user_id, from_user_id,
					target_type, target_detailed_type, 
					target_id, target_uuid,
					message, data) 
					VALUES (?, ?, ?, ?, ?, ?, ?, ?);`
	result, insertErr := sqlDB.Exec(insertQuery,
		notification.ToUserId, notification.FromUserId,
		notification.TargetType, notification.TargetDetailedType,
		notification.TargetId, notification.TargetUUID,
		notification.Message, notification.Data)
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

func InsertNotificationForEvent(n *Notification, groupID, userID int) error {
	qry := `INSERT INTO notifications (
				to_user_id, from_user_id,
				target_type, target_detailed_type,
				target_id, target_uuid,
				message
			)
			SELECT 
				members.follower_id, ?,
				'groups', 'group_event',
				?, ?, ?
			FROM following members
			WHERE members.group_id = ? AND members.follower_id != ?
				AND members.status = 'accepted';`

	_, err := sqlDB.Exec(qry,
		n.FromUserId,
		n.TargetId, n.TargetUUID, n.Message, groupID, userID)
	return err
}

func UpdateNotificationOnCancel(n *Notification, tgtType, reqType string) error {
	qry := `WITH target AS (
				SELECT id
				FROM notifications
				WHERE to_user_id = ? AND from_user_id = ?
					AND target_type = '` + tgtType + `' AND message = '` + reqType + `'
				ORDER BY created_at DESC
				LIMIT 1
			)
			UPDATE notifications
			SET message = ?, updated_by = ?, updated_at = CURRENT_TIMESTAMP,
				status = 'delete'
			WHERE id IN (SELECT id FROM target);`

	_, err := sqlDB.Exec(qry, n.ToUserId, n.FromUserId, n.Message, n.UpdatedBy)
	return err
}

func UpdateNotificationReadStatus(notification_id int, is_read int, user_id int) error {
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
        SELECT 
			n.id AS notification_id,
			n.to_user_id AS notification_to_user_id,
			n.from_user_id AS notification_from_user_id,
			n.target_id AS notification_target_id,
			n.target_uuid AS notification_target_uuid,
			n.target_type AS notification_target_type,
			n.target_detailed_type AS notification_target_detailed_type,

			IFNULL(
				CASE
					WHEN n.target_detailed_type = 'follow_request' THEN
						'You have a follow request from <b>' || target_user.nick_name || '</b>'
					WHEN n.target_detailed_type = 'follow_request_responded' THEN
						CASE
							WHEN n.message = 'You have a new follower' THEN
								'You have a new follower: <b>' || target_user.nick_name || '</b>'
							ELSE
								'Your follow request to <b>' || target_user.nick_name || '</b> has been ' || n.message
						END
					WHEN n.target_detailed_type = 'group_invite' THEN
						'You have been invited to group <b>' || target_group.title || '</b>'
					WHEN n.target_detailed_type = 'group_invite_responded' THEN
						'Your group invitation to group <b>' || target_group.title || '</b> has been ' || n.message
					WHEN n.target_detailed_type = 'group_request' THEN
						'You have a group joining request from <b>' || from_user.nick_name || '</b> to group <b>' || target_group.title || '</b>'
					WHEN n.target_detailed_type = 'group_request_responded' THEN
						'Your group joining request to group <b>' || target_group.title || '</b> has been ' || n.message
					WHEN n.target_detailed_type = 'group_event' THEN
						'Event <b>' || n.message || '</b>' || ' has been created in group <b>' || target_group.title || '</b>'
					ELSE ''
				END,
				''
			) AS notification_message,

			n.is_read AS notification_is_read,
			n.data AS notification_data,
			n.status AS notification_status,
			n.created_at AS notification_created_at,
			n.updated_at AS notification_updated_at,
			n.updated_by AS notification_updated_by,

			from_user.id AS from_user_id,
			from_user.uuid AS from_user_uuid,
			from_user.first_name AS from_user_first_name,
			from_user.last_name AS from_user_last_name,
			from_user.nick_name AS from_user_nick_name,
			from_user.profile_image AS from_user_profile_image,

			to_user.uuid AS to_user_uuid

			FROM notifications n
			INNER JOIN users to_user
				ON n.to_user_id = to_user.id
				AND to_user.id = ?
			LEFT JOIN users from_user
				ON n.from_user_id = from_user.id
			LEFT JOIN users target_user
				ON n.target_detailed_type IN ('follow_request', 'follow_request_responded')
				AND n.target_id = target_user.id
			LEFT JOIN groups target_group
				ON n.target_detailed_type IN ('group_invite', 'group_invite_responded', 'group_request', 'group_request_responded', 'group_event')
				AND n.target_id = target_group.id

			WHERE n.status != 'delete'
			ORDER BY n.id DESC;
    `, to_user_id)
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	var notifications []Notification

	for rows.Next() {
		var notification Notification
		var fromUser NotificationFromUser

		// Scan the post and user data
		err := rows.Scan(
			&notification.ID, &notification.ToUserId, &notification.FromUserId,
			&notification.TargetId, &notification.TargetUUIDForm, &notification.TargetType, &notification.TargetDetailedType,
			&notification.Message, &notification.IsRead, &notification.Data,
			&notification.Status, &notification.CreatedAt, &notification.UpdatedAt, &notification.UpdatedBy,
			&fromUser.ID, &fromUser.UUID, &fromUser.FirstName, &fromUser.LastName, &fromUser.NickName, &fromUser.ProfileImage,
			&notification.ToUserUUID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

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
        SELECT 
			n.id AS notification_id,
			n.to_user_id AS notification_to_user_id,
			n.from_user_id AS notification_from_user_id,
			n.target_id AS notification_target_id,
			n.target_uuid AS notification_target_uuid,
			n.target_type AS notification_target_type,
			n.target_detailed_type AS notification_target_detailed_type,

			IFNULL(
				CASE
					WHEN n.target_detailed_type = 'follow_request' THEN
						'You have a follow request from <b>' || target_user.nick_name || '</b>'
					WHEN n.target_detailed_type = 'follow_request_responded' THEN
						CASE
							WHEN n.message = 'You have a new follower' THEN
								n.message || ': <b>' || target_user.nick_name || '</b>'
							ELSE
								'Your follow request to <b>' || target_user.nick_name || '</b> has been ' || n.message
						END
					WHEN n.target_detailed_type = 'group_invite' THEN
						'You have been invited to group <b>' || target_group.title || '</b>'
					WHEN n.target_detailed_type = 'group_invite_responded' THEN
						'Your group invitation to group <b>' || target_group.title || '</b> has been ' || n.message
					WHEN n.target_detailed_type = 'group_request' THEN
						'You have a group joining request from <b>' || from_user.nick_name || '</b> to group <b>' || target_group.title || '</b>'
					WHEN n.target_detailed_type = 'group_request_responded' THEN
						'Your group joining request to group <b>' || target_group.title || '</b> has been ' || n.message
					WHEN n.target_detailed_type = 'group_event' THEN
						'Event <b>' || n.message || '</b>' || ' has been created in group <b>' || target_group.title || '</b>'
					ELSE ''
				END,
				''
			) AS notification_message,

			n.is_read AS notification_is_read,
			n.data AS notification_data,
			n.status AS notification_status,
			n.created_at AS notification_created_at,
			n.updated_at AS notification_updated_at,
			n.updated_by AS notification_updated_by,

			from_user.id AS from_user_id,
			from_user.uuid AS from_user_uuid,
			from_user.first_name AS from_user_first_name,
			from_user.last_name AS from_user_last_name,
			from_user.nick_name AS from_user_nick_name,
			from_user.profile_image AS from_user_profile_image,

			to_user.uuid AS to_user_uuid

			FROM notifications n
			INNER JOIN users to_user
				ON n.to_user_id = to_user.id
				AND to_user.id = ?
			LEFT JOIN users from_user
				ON n.from_user_id = from_user.id
			LEFT JOIN users target_user
				ON n.target_detailed_type IN ('follow_request', 'follow_request_responded')
				AND n.target_id = target_user.id
			LEFT JOIN groups target_group
				ON n.target_detailed_type IN ('group_invite', 'group_invite_responded', 'group_request', 'group_request_responded', 'group_event')
				AND n.target_id = target_group.id

			WHERE n.status != 'delete'
				AND n.id = ?
			ORDER BY n.id DESC;
    `, to_user_id, notification_id)
	if selectError != nil {
		return Notification{}, selectError
	}
	defer rows.Close()

	var notification Notification

	if rows.Next() {
		var fromUser NotificationFromUser

		// Scan the post and user data
		err := rows.Scan(
			&notification.ID, &notification.ToUserId, &notification.FromUserId,
			&notification.TargetId, &notification.TargetUUIDForm, &notification.TargetType, &notification.TargetDetailedType,
			&notification.Message, &notification.IsRead, &notification.Data,
			&notification.Status, &notification.CreatedAt, &notification.UpdatedAt, &notification.UpdatedBy,
			&fromUser.ID, &fromUser.UUID, &fromUser.FirstName, &fromUser.LastName, &fromUser.NickName, &fromUser.ProfileImage,
			&notification.ToUserUUID,
		)
		if err != nil {
			return Notification{}, fmt.Errorf("error scanning row: %v", err)
		}

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
