package dbTools

import "strconv"

func (db *DBContainer) InsertMessage(m *Message) error {
	qry := `INSERT INTO messages 
			(sender_id, receiver_id, content) 
			VALUES (?, ?, ?)`
	_, err := db.conn.Exec(qry,
		m.SenderID,
		m.ReceiverID,
		m.Content,
	)
	return err
}

func (db *DBContainer) UpdateMessage(m *Message) error {
	qry := `UPDATE messages SET read_at = ? WHERE id = ?`
	_, err := db.conn.Exec(qry,
		m.ReadAt,
		m.ID,
	)
	return err
}

func (db *DBContainer) SelectMessages(id_1, id_2 int, msgIdStr string) (*[]Message, error) {
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

	rows, err := db.conn.Query(qry, id_1, id_2, id_2, id_1, msgId)
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

func (db *DBContainer) SelectUnreadMessages(senderID, receiverID int) (*[]Message, error) {
	qry := `SELECT * FROM messages
			WHERE (sender_id = ? AND receiver_id = ? AND read_at IS NULL)`

	rows, err := db.conn.Query(qry, senderID, receiverID)
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

func (db *DBContainer) SelectUserList(receiverID int) (*[]string, *[]int, error) {
	qry := `SELECT u.nick_name, u.id
			FROM users U
			LEFT JOIN (
				SELECT sender_id, receiver_id, created_at
				FROM messages
				WHERE (receiver_id = ? OR sender_id = ?)
			) m ON (u.id = m.receiver_id AND m.sender_id = ? OR u.id = m.sender_id AND m.receiver_id  = ?)
			WHERE u.id != 0
			GROUP BY u.nick_name, u.id
			ORDER BY MAX(m.created_at) DESC, LOWER(u.nick_name) ASC`

	rows, err := db.conn.Query(qry, receiverID, receiverID, receiverID, receiverID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var names []string
	var ids []int
	for rows.Next() {
		var n string
		var id int
		err := rows.Scan(&n, &id)
		if err != nil {
			return nil, nil, err
		}
		names = append(names, n)
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, checkErrNoRows(err)
	}
	return &names, &ids, nil
}

func (db *DBContainer) SelectUnreadMsgList(receiverID int) (*[]int, error) {
	qry := `SELECT DISTINCT u.id
			FROM messages m
			JOIN users u ON m.sender_id = u.id
			WHERE (m.receiver_id = ?) AND read_at IS NULL`

	rows, err := db.conn.Query(qry, receiverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, checkErrNoRows(err)
	}
	return &ids, nil
}
