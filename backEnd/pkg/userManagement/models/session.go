package models

import (
	"database/sql"
	"fmt"
	"social-network/pkg/utils"
	"time"
)

// User struct represents the user data model
type Session struct {
	ID         string    `json:"id"`
	UserId     int       `json:"user_id"`
	IsActive   bool      `json:"is_active"`
	StartTime  time.Time `json:"start_time"`
	ExpireTime time.Time `json:"expire_time"`
	LastAccess time.Time `json:"last_access"`
}

var sqlDB *sql.DB

func Initialize(dbMain *sql.DB) {
	sqlDB = dbMain
}

// checkErrNoRows() checks if no result from sql query.
func checkErrNoRows(err error) error {
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

// to check and remove expired sessions
func SelectActiveSessionBy(field string, id interface{}) (*Session, error) {
	if field != "id" && field != "user_id" {
		return nil, fmt.Errorf("invalid field")
	}
	var s Session
	qry := `SELECT * FROM sessions WHERE ` + field + ` = ? AND is_active = 1`
	err := sqlDB.QueryRow(qry, id).Scan(
		&s.ID,
		&s.UserId,
		&s.IsActive,
		&s.StartTime,
		&s.ExpireTime,
		&s.LastAccess)
	return &s, err
}

// to do update then insert?
func InsertSession(session *Session) (*Session, error) {
	// Generate UUID for the user if not already set
	sessionId, err := utils.GenerateUuid()
	if err != nil {
		return nil, err
	}
	session.ID = sessionId

	// Set session expiration time to 1 hour
	session.ExpireTime = time.Now().Add(1 * time.Hour)

	qry := `INSERT INTO sessions
			(id, user_id, is_active, expire_time)
			VALUES ( ?, ?, ?, ?)`
	_, err = sqlDB.Exec(qry,
		session.ID,
		session.UserId,
		session.IsActive,
		session.ExpireTime)

	return session, err
}

// UpdateSession() for when Session is expired, logout or refreshed
func UpdateSession(s *Session) error {
	qry := `UPDATE sessions
			SET is_active = ?, expire_time = ?, last_access= ?
			WHERE id = ?`
	_, err := sqlDB.Exec(qry,
		s.IsActive,
		s.ExpireTime,
		s.LastAccess,
		s.ID)
	return err
}

/* func SelectSession(sessionToken string) (User, time.Time, error) {
	var user User
	var expirationTime time.Time
	// todo: fix type_id, profile_image (, IFNULL(u.profile_image, '') as profile_image)
	err := sqlDB.QueryRow(`SELECT
							u.id as user_id, u.type_id as user_type_id, u.first_name as user_first_name, u.last_name as user_last_name, u.gender as user_gender, u.birthday as user_birthday, u.nick_name as nick_name, u.email as user_email, IFNULL(u.profile_image, '') as profile_image,
							expire_time
						FROM sessions s
							INNER JOIN users u
								ON s.user_id = u.id
						WHERE s.id = ?`, sessionToken).Scan(&user.ID, &user.TypeId, &user.FirstName, &user.LastName, &user.Gender, &user.BirthDay, &user.NickName, &user.Email, &user.ProfileImage, &expirationTime)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			// Handle other database errors
			return User{}, time.Time{}, errors.New("sql: no rows in result set")
		} else {
			fmt.Println("error is:", err)
			// Handle other database errors
			return User{}, time.Time{}, errors.New("database error")
		}
	}

	return user, expirationTime, nil
}

// IsSessionActive checks if a session is active based on the session token
func IsSessionActive(sessionToken string) (bool, error) {
	var expiresAt time.Time

	// Query the database for the session's expiration time
	err := sqlDB.QueryRow(`SELECT expire_time FROM sessions WHERE id = ?`, sessionToken).Scan(&expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// No session found for the given token
			return false, errors.New("session not found")
		}
		// Handle other database errors
		return false, err
	}

	// Check if the session is still active
	if expiresAt.After(time.Now()) {
		return true, nil // Session is active
	}

	return false, nil // Session is expired
}
*/
