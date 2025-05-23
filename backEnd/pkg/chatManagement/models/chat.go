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
	Content      string       `json:"content"`
	CreatedAt    time.Time    `json:"createdAt"`
	ReadAt       sql.NullTime `json:"readAt"`
}

func InsertMessage(m *Message) error {
	qry := `INSERT INTO messages 
			(sender_id, receiver_id, content) 
			VALUES (?, ?, ?)`
	_, err := sqlDB.Exec(qry,
		m.SenderID,
		m.ReceiverID,
		m.Content,
	)
	return err
}

func UpdateMessage(m *Message) error {
	qry := `UPDATE messages SET read_at = ? WHERE id = ?`
	_, err := sqlDB.Exec(qry,
		m.ReadAt,
		m.ID,
	)
	return err
}

func SelectMessages(id_1, id_2 int, msgIdStr string) (*[]Message, error) {
	msgId, err := strconv.Atoi(msgIdStr)
	if err != nil {
		return nil, err
	}

	fromStart := ""
	if msgId != -1 {
		fromStart = "AND id < ?"
	}
	qry := `SELECT * FROM messages
			WHERE (sender_id = ? AND receiver_id = ? OR sender_id = ? AND receiver_id = ?)
			` + fromStart + `
			ORDER BY created_at DESC
			LIMIT 10`

	rows, err := sqlDB.Query(qry, id_1, id_2, id_2, id_1, msgId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		err := rows.Scan(
			&m.ID,
			&m.SenderID,
			&m.ReceiverID,
			&m.Content,
			&m.CreatedAt,
			&m.ReadAt,
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

func SelectUnreadMessages(senderID, receiverID int) (*[]Message, error) {
	qry := `SELECT * FROM messages
			WHERE (sender_id = ? AND receiver_id = ? AND read_at IS NULL)`

	rows, err := sqlDB.Query(qry, senderID, receiverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		err := rows.Scan(
			&m.ID,
			&m.SenderID,
			&m.ReceiverID,
			&m.Content,
			&m.CreatedAt,
			&m.ReadAt,
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
			LEFT JOIN (
				SELECT sender_id, receiver_id, created_at
				FROM messages
				WHERE (receiver_id = ? OR sender_id = ?)
			) m ON (u.id = m.receiver_id AND m.sender_id = ? OR u.id = m.sender_id AND m.receiver_id  = ?)
			WHERE u.id != 0
			GROUP BY u.nick_name, u.uuid
			ORDER BY MAX(m.created_at) DESC, LOWER(u.nick_name) ASC`

	rows, err := sqlDB.Query(qry, receiverID, receiverID, receiverID, receiverID)
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
			WHERE (m.receiver_id = ?) AND read_at IS NULL`

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
