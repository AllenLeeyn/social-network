package model

import (
	"database/sql"
	"strconv"
	"time"
)

type Message struct {
	ID           int    `json:"ID"`
	Action       string `json:"action"`
	SenderUUID   string `json:"senderUUID"`
	SenderID     int
	ReceiverUUID string `json:"receiverUUID"`
	ReceiverID   int
	GroupUUID    string `json:"groupUUID"`
	GroupID      int
	Content      string         `json:"content"`
	Status       sql.NullString `json:"status"` // new
	CreatedAt    time.Time      `json:"createdAt"`
	ReadAt       sql.NullTime   `json:"readAt"`
	UpdatedBy    sql.NullInt64  `json:"updated_by"` // new
	UpdatedAt    sql.NullTime   `json:"updated_at"` // new
}

type MessageView struct {
	ID           int          `json:"ID"`
	Action       string       `json:"action"`
	SenderUUID   string       `json:"senderUUID"`
	SenderName   string       `json:"senderName"`
	ReceiverUUID string       `json:"receiverUUID"`
	ReceiverName string       `json:"ReceiverName"`
	GroupUUID    string       `json:"groupUUID"`
	Content      string       `json:"content"`
	CreatedAt    time.Time    `json:"createdAt"`
	ReadAt       sql.NullTime `json:"readAt"`
}

func InsertMessage(m *Message) (int, error) {
	qry := `INSERT INTO messages (sender_id, receiver_id, content, group_id)
			VALUES (
			(SELECT id FROM users WHERE uuid = ?),
			(SELECT id FROM users WHERE uuid = ?),
			?,
			(SELECT id FROM groups WHERE uuid = ?)
			)`
	result, err := sqlDB.Exec(qry,
		m.SenderUUID,
		m.ReceiverUUID,
		m.Content,
		m.GroupUUID,
	)
	if err != nil {
		return -1, err
	}

	msgID, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(msgID), err
}

func UpdateMessage(m *Message) error {
	qry := `UPDATE messages SET read_at = ? WHERE id = ?`
	_, err := sqlDB.Exec(qry,
		m.ReadAt,
		m.ID,
	)
	return err
}

func SelectMessage(msgID int) (*MessageView, error) {

	qry := `SELECT  m.id, sender.uuid, sender.nick_name,
					receiver.uuid, receiver.nick_name, g.uuid,
					m.content, m.created_at, m.read_at
			FROM messages m
			JOIN users sender ON m.sender_id = sender.id
			JOIN users receiver ON m.receiver_id = receiver.id
			JOIN groups g ON m.group_id = g.id
			WHERE m.id = ?
			LIMIT 1;`

	row := sqlDB.QueryRow(qry, msgID)

	var m MessageView
	err := row.Scan(
		&m.ID, &m.SenderUUID, &m.SenderName,
		&m.ReceiverUUID, &m.ReceiverName, &m.GroupUUID,
		&m.Content, &m.CreatedAt, &m.ReadAt,
	)
	m.Action = "message"
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func SelectPrivateMessages(uuid_1, uuid_2, msgIdStr string) (*[]MessageView, error) {
	msgId, err := strconv.Atoi(msgIdStr)
	if err != nil {
		return nil, err
	}

	qry := `SELECT  m.id, sender.uuid, sender.nick_name,
					receiver.uuid, receiver.nick_name,
					'00000000-0000-0000-0000-000000000000',
					m.content, m.created_at, m.read_at
			FROM messages m
			JOIN users sender ON m.sender_id = sender.id
			JOIN users receiver ON m.receiver_id = receiver.id
			WHERE 
				(
					(sender.uuid = ? AND receiver.uuid = ?)
					OR
					(sender.uuid = ? AND receiver.uuid = ?)
				)
				AND m.group_id = 0
				AND m.status = 'enable'
				AND (m.id < ? OR ? = -1)
			ORDER BY m.created_at DESC
			LIMIT 10;`

	rows, err := sqlDB.Query(qry, uuid_1, uuid_2, uuid_2, uuid_1, msgId, msgId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []MessageView
	for rows.Next() {
		var m MessageView
		err := rows.Scan(
			&m.ID, &m.SenderUUID, &m.SenderName,
			&m.ReceiverUUID, &m.ReceiverName, &m.GroupUUID,
			&m.Content, &m.CreatedAt, &m.ReadAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	if err := rows.Err(); err != nil {
		return nil, checkErrNoRows(err)
	}
	return &messages, nil
}

func SelectGroupMessages(_, groupUUID, msgIdStr string) (*[]MessageView, error) {
	msgId, err := strconv.Atoi(msgIdStr)
	if err != nil {
		return nil, err
	}

	qry := `SELECT  m.id, sender.uuid, sender.nick_name,
					receiver.uuid, receiver.nick_name, g.uuid,
					m.content, m.created_at, m.read_at
			FROM messages m
			JOIN users sender ON m.sender_id = sender.id
			JOIN users receiver ON m.receiver_id = receiver.id
			JOIN groups g ON m.group_id = g.id
			WHERE g.uuid = ?
			  AND m.status = 'enable'
			  AND (m.id < ? OR ? = -1)
			ORDER BY m.created_at DESC
			LIMIT 10;`

	rows, err := sqlDB.Query(qry, groupUUID, msgId, msgId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []MessageView
	for rows.Next() {
		var m MessageView
		err := rows.Scan(
			&m.ID, &m.SenderUUID, &m.SenderName,
			&m.ReceiverUUID, &m.ReceiverName, &m.GroupUUID,
			&m.Content, &m.CreatedAt, &m.ReadAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	if err := rows.Err(); err != nil {
		return nil, checkErrNoRows(err)
	}
	return &messages, nil
}

func SelectUnreadMessages(senderUUID, receiverUUID string) (*[]Message, error) {
	qry := `SELECT m.id, sender.uuid,
				   receiver.uuid,
				   '00000000-0000-0000-0000-000000000000',
				   m.content, m.created_at, m.read_at
			FROM messages m
			JOIN users sender ON m.sender_id = sender.id
			JOIN users receiver ON m.receiver_id = receiver.id
			WHERE m.group_id = 0
			  AND sender.uuid = ?
			  AND receiver.uuid = ?
			  AND m.read_at IS NULL
			  AND m.status = 'enable'`

	rows, err := sqlDB.Query(qry, senderUUID, senderUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		err := rows.Scan(
			&m.ID, &m.SenderUUID, &m.ReceiverUUID, &m.GroupUUID,
			&m.Content, &m.CreatedAt, &m.ReadAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	if err := rows.Err(); err != nil {
		return nil, checkErrNoRows(err)
	}
	return &messages, nil
}

func SelectUserList(receiverID int) (*[]string, *[]string, error) {
	qry := `SELECT u.nick_name, u.uuid
			FROM users U
			JOIN following f ON (
				(f.follower_id = ? AND f.leader_id = u.id)
				OR
				(f.leader_id = ? AND f.follower_id = u.id))
				AND f.group_id = 0
				AND f.status = 'accepted'
			LEFT JOIN (
  				SELECT 
					CASE 
						WHEN sender_id = ? THEN receiver_id 
						ELSE sender_id 
					END AS user_id, MAX(created_at) AS created_at
				FROM messages
				WHERE (receiver_id = ? OR sender_id = ?)
      				AND group_id = 0
    				AND status = 'enable'
   				GROUP BY user_id
			) m ON m.user_id = u.id
			WHERE u.id != 0
			ORDER BY m.created_at DESC, LOWER(u.nick_name) ASC`

	rows, err := sqlDB.Query(qry,
		receiverID, receiverID, receiverID, receiverID, receiverID, receiverID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var names []string
	var uuids []string
	for rows.Next() {
		var n string
		var uuid string
		err := rows.Scan(&n, &uuid)
		if err != nil {
			return nil, nil, err
		}
		names = append(names, n)
		uuids = append(uuids, uuid)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, checkErrNoRows(err)
	}
	return &names, &uuids, nil
}

func SelectUnreadMsgList(receiverID int) (*[]string, error) {
	qry := `SELECT DISTINCT u.uuid
			FROM messages m
			JOIN users u ON m.sender_id = u.id
			WHERE (m.receiver_id = ?) 
				AND m.read_at IS NULL
				AND m.group_id = 0`

	rows, err := sqlDB.Query(qry, receiverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uuids []string
	for rows.Next() {
		var uuid string
		err := rows.Scan(&uuid)
		if err != nil {
			return nil, err
		}
		uuids = append(uuids, uuid)
	}
	if err := rows.Err(); err != nil {
		return nil, checkErrNoRows(err)
	}
	return &uuids, nil
}
