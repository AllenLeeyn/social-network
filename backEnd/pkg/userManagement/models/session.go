package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
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

func InsertSession(sqlDB *sql.DB, session *Session) (*Session, error) {
	// Generate UUID for the user if not already set
	if session.ID == "" {
		uuidSessionTokenid, err := utils.GenerateUuid()
		if err != nil {
			return nil, err
		}
		session.ID = uuidSessionTokenid
	}

	// Set session expiration time to 1 hour
	session.ExpireTime = time.Now().Add(1 * time.Hour)

	// Start a transaction for atomicity
	tx, err := sqlDB.Begin()
	if err != nil {
		return &Session{}, err
	}

	updateQuery := `UPDATE sessions SET expire_time = CURRENT_TIMESTAMP WHERE user_id = ? AND expire_time > CURRENT_TIMESTAMP;`
	_, updateErr := tx.Exec(updateQuery, session.UserId)
	if updateErr != nil {
		tx.Rollback()
		return nil, updateErr
	}

	insertQuery := `INSERT INTO sessions (id, user_id, expire_time, is_active) VALUES (?, ?, ?, ?);`
	_, insertErr := tx.Exec(insertQuery, session.ID, session.UserId, session.ExpireTime, session.IsActive)
	if insertErr != nil {
		tx.Rollback()
		// Check if the error is a SQLite constraint violation
		if sqliteErr, ok := insertErr.(interface{ ErrorCode() int }); ok {
			if sqliteErr.ErrorCode() == 19 { // SQLite constraint violation error code
				return nil, sql.ErrNoRows // Return custom error to indicate a duplicate
			}
		}
		return nil, insertErr
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		tx.Rollback() // Rollback on error
		return nil, err
	}

	return session, nil
}

func SelectSession(sqlDB *sql.DB, sessionToken string) (User, time.Time, error) {
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

func DeleteSession(sqlDB *sql.DB, sessionToken string) error {
	_, err := sqlDB.Exec(`UPDATE sessions
					SET expire_time = CURRENT_TIMESTAMP
					WHERE id = ?;`, sessionToken)
	if err != nil {
		// Handle other database errors
		log.Fatal(err)
		return errors.New("database error")
	}

	return nil

}

// IsSessionActive checks if a session is active based on the session token
func IsSessionActive(sqlDB *sql.DB, sessionToken string) (bool, error) {
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

func GetUserIDFromCookie(sqlDB *sql.DB, r *http.Request) (int, string, error) {
	// Retrieve user data (e.g., from session or database)
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		// Return an error if the session token is not found
		return 0, "", fmt.Errorf("error retrieving session token: %v", err)
	}

	// Fetch the user from the database using the session token
	user, _, err := SelectSession(sqlDB, sessionToken.Value)
	if err != nil {
		return 0, "", fmt.Errorf("error retrieving session: %v", err)
	}

	myUserID := user.ID         // Get the ID field from the User struct
	myUsername := user.NickName // Use the Username field from the User struct

	return myUserID, myUsername, nil
}
