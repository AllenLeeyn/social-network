package dbTools

import (
	"fmt"
)

/*
	db.SelectActiveSessionBy() id or user_id.

valid field: "user_id", "id".

"user_id" is for checking if user has active session when logging in.
"id" is for finding a user_id based on the current session.
*/
func (db *DBContainer) SelectActiveSessionBy(field string, id interface{}) (*Session, error) {
	if field != "id" && field != "user_id" {
		return nil, fmt.Errorf("invalid field")
	}
	var s Session
	qry := `SELECT * FROM sessions WHERE ` + field + ` = ? AND is_active = 1`
	err := db.Conn.QueryRow(qry, id).Scan(
		&s.ID,
		&s.UserID,
		&s.IsActive,
		&s.StartTime,
		&s.ExpireTime,
		&s.LastAccess)
	return &s, err
}

func (db *DBContainer) SelectActiveSessions() (*[]Session, error) {
	qry := `SELECT * FROM sessions WHERE is_active = 1`

	rows, err := db.Conn.Query(qry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var s Session
		err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.IsActive,
			&s.StartTime,
			&s.ExpireTime,
			&s.LastAccess)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	if err := rows.Err(); err != nil {
		return nil, checkErrNoRows(err)
	}
	return &sessions, nil
}

// db.InsertSession() when User login is successful
func (db *DBContainer) InsertSession(s *Session) error {
	qry := `INSERT INTO sessions
			(id, user_id, is_active, expire_time)
			VALUES ( ?, ?, ?, ?)`
	_, err := db.Conn.Exec(qry,
		s.ID,
		s.UserID,
		s.IsActive,
		s.ExpireTime)
	return err
}

// db.UpdateSession() for when Session is expired, logout or refreshed
func (db *DBContainer) UpdateSession(s *Session) error {
	qry := `UPDATE sessions
			SET is_active = ?, expire_time = ?, last_access= ?
			WHERE id = ?`
	_, err := db.Conn.Exec(qry,
		s.IsActive,
		s.ExpireTime,
		s.LastAccess,
		s.ID)
	return err
}
