package models

import (
	"fmt"
	"social-network/pkg/utils"
	"time"
)

type Session struct {
	ID         string    `json:"id"`
	UserId     int       `json:"user_id"`
	IsActive   bool      `json:"is_active"`
	StartTime  time.Time `json:"start_time"`
	ExpireTime time.Time `json:"expire_time"`
	LastAccess time.Time `json:"last_access"`
}

const sessionDuration = time.Hour

func (um *UserModel) SelectActiveSessionBy(field string, id interface{}) (*Session, error) {
	if field != "id" && field != "user_id" {
		return nil, fmt.Errorf("invalid field")
	}
	var s Session
	qry := `SELECT * FROM sessions 
			WHERE ` + field + ` = ? AND is_active = 1 AND expire_time > CURRENT_TIMESTAMP`
	err := um.DB.QueryRow(qry, id).Scan(
		&s.ID,
		&s.UserId,
		&s.IsActive,
		&s.StartTime,
		&s.ExpireTime,
		&s.LastAccess)
	return &s, err
}

func (um *UserModel) InsertSession(session *Session) (*Session, error) {
	sessionId, err := utils.GenerateUuid()
	if err != nil {
		return nil, err
	}
	session.ID = sessionId

	// Set session expiration time to 1 hour
	session.ExpireTime = time.Now().Add(sessionDuration)

	qry := `INSERT INTO sessions
			(id, user_id, is_active, expire_time)
			VALUES ( ?, ?, ?, ?)`
	_, err = um.DB.Exec(qry,
		session.ID,
		session.UserId,
		session.IsActive,
		session.ExpireTime)

	return session, err
}

// UpdateSession() for when Session is expired, logout or refreshed
func (um *UserModel) UpdateSession(s *Session) error {
	qry := `UPDATE sessions
			SET is_active = ?, expire_time = ?, last_access= ?
			WHERE id = ?`
	_, err := um.DB.Exec(qry,
		s.IsActive,
		s.ExpireTime,
		s.LastAccess,
		s.ID)
	return err
}
