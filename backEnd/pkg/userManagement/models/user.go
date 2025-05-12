package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"social-network/pkg/dbTools"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User struct represents the user data model
type User struct {
	ID           int        `json:"id"`
	UUID         string     `json:"uuid"`
	TypeId       int        `json:"type_id"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Gender       string     `json:"gender"`
	BirthDay     time.Time  `json:"birthday"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"pw_hash"`
	NickName     string     `json:"nick_name"`
	ProfileImage string     `json:"profile_image"`
	AboutMe      string     `json:"about_me"`
	Visibility   string     `json:"visibility"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedBy    int        `json:"updated_by"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

func InsertUser(db *dbTools.DBContainer, user *User) (int, error) {
	var existingEmail string
	var existingUsername string
	emailCheckQuery := `SELECT email, nick_name FROM users WHERE email = ? OR nick_name = ? LIMIT 1;`
	err := db.Conn.QueryRow(emailCheckQuery, user.Email, user.NickName).Scan(&existingEmail, &existingUsername)
	if err == nil {
		if existingEmail == user.Email {
			return -1, errors.New("duplicateEmail")
		}
		if existingUsername == user.NickName {
			return -1, errors.New("duplicateUsername")
		}
	}

	// todo: fix type_id
	insertQuery := `INSERT INTO users (first_name, last_name, type_id, birthday, gender, nick_name, email, pw_hash, about_me, visibility) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	result, insertErr := db.Conn.Exec(insertQuery, user.FirstName, user.LastName, 1, user.BirthDay, user.Gender, user.NickName, user.Email, user.PasswordHash, user.AboutMe, user.Visibility)
	if insertErr != nil {
		// Check if the error is a SQLite constraint violation (duplicate entry)
		if sqliteErr, ok := insertErr.(interface{ ErrorCode() int }); ok {
			if sqliteErr.ErrorCode() == 19 { // 19 = UNIQUE constraint failed (SQLite error code)
				return -1, errors.New("user with this email or nick_name already exists")
			}
		}
		return -1, insertErr // Other DB errors
	}

	// Retrieve the last inserted ID
	userId, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return -1, err
	}

	return int(userId), nil
}

func UpdateUser(db *dbTools.DBContainer, user *User) error {
	if user.ProfileImage == "" {
		updateUser := `UPDATE users
					SET first_name = ?,
						last_name = ?,
						gender = ?,
						birthday = ?,
						about_me = ?,
						visibility = ?,
						updated_at = CURRENT_TIMESTAMP,
						updated_by = ?
					WHERE id = ?;`
		_, updateErr := db.Conn.Exec(updateUser, user.FirstName, user.LastName, user.Gender, user.BirthDay, user.AboutMe, user.Visibility, user.ID, user.ID)

		if updateErr != nil {
			return updateErr
		}
	} else {
		updateUser := `UPDATE users
		SET first_name = ?,
			last_name = ?,
			gender = ?,
			birthday = ?,
			about_me = ?,
			visibility = ?,
			profile_image = ?,
			updated_at = CURRENT_TIMESTAMP,
			updated_by = ?
		WHERE id = ?;`
		_, updateErr := db.Conn.Exec(updateUser, user.FirstName, user.LastName, user.Gender, user.BirthDay, user.AboutMe, user.Visibility, user.ProfileImage, user.ID, user.ID)

		if updateErr != nil {
			return updateErr
		}
	}
	return nil
}

func AuthenticateUser(db *dbTools.DBContainer, nick_name, password string) (bool, int, error) {
	// Query to retrieve the hashed password stored in the database for the given nick_name
	var userId int
	var storedHashedPassword string
	err := db.Conn.QueryRow("SELECT id, pw_hash FROM users WHERE nick_name = ? or email = ?", nick_name, nick_name).Scan(&userId, &storedHashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			// NickName not found
			return false, -1, errors.New("nick_name not found")
		}
		// Handle other database errors
		log.Fatal(err)
	}

	// Compare the entered password with the stored hashed password using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(password))
	if err != nil {
		// PasswordHash is incorrect
		return false, -1, errors.New("password is incorrect")
	}

	// Successful login if no errors occurred
	return true, userId, nil
}

func ReadAllUsers(db *dbTools.DBContainer) ([]User, error) {
	// Query the records
	rows, selectError := db.Conn.Query(`
        SELECT u.id as user_id, u.first_name as user_first_name, u.last_name as user_last_name, u.nick_name as user_nick_name, u.email as user_email, u.about_me as user_about_me, u.visibility as user_visibility,
		IFNULL(u.profile_image, '') as profile_image, u.status as user_status, u.created_at as user_created_at, 
		u.updated_at as user_updated_at, u.updated_by as user_updated_by
		FROM users u
		WHERE u.status != 'delete'
		AND u.type != 'admin'
		ORDER BY u.id desc;
    `)
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		// Scan the post and user data
		err := rows.Scan(
			&user.ID, &user.FirstName, &user.LastName, &user.NickName, &user.Email, &user.AboutMe, &user.Visibility,
			&user.ProfileImage, &user.Status, &user.CreatedAt,
			&user.UpdatedAt, &user.UpdatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		users = append(users, user)
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return users, nil
}

func ReadAllChatUsers(db *dbTools.DBContainer, user_id int) ([]User, error) {
	rows, selectError := db.Conn.Query(`
        SELECT u.id as user_id, u.first_name as user_first_name, u.last_name as user_last_name, u.nick_name as user_nick_name, u.email as user_email, u.about_me as user_about_me, u.visibility as user_visibility, 
		IFNULL(u.profile_image, '') as profile_image, u.status as user_status, u.created_at as user_created_at, 
		u.updated_at as user_updated_at, u.updated_by as user_updated_by
			FROM users u
			LEFT JOIN chat_members cm 
				ON u.id = cm.user_id
				AND cm.status != 'delete'
			LEFT JOIN chats c
				ON cm.chat_id = c.id
				AND c.status != 'delete'
				AND c.id in (
					SELECT chat_id
					FROM chat_members
					WHERE user_id = ?
					AND status != 'delete'
				)
			WHERE u.status != 'delete'
			AND u.type != 'admin'
			GROUP BY u.id
			ORDER BY MAX(c.updated_at) desc, u.nick_name;
    `, user_id)
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.ID, &user.FirstName, &user.LastName, &user.NickName, &user.Email, &user.AboutMe, &user.Visibility,
			&user.ProfileImage, &user.Status, &user.CreatedAt,
			&user.UpdatedAt, &user.UpdatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return users, nil
}

func ReadUserByID(db *dbTools.DBContainer, user_id int) (User, error) {
	// Query the records
	rows, selectError := db.Conn.Query(`
        SELECT u.id as user_id, u.first_name as user_first_name, u.last_name as user_last_name, u.nick_name as user_nick_name, u.email as user_email, u.about_me as user_about_me, u.visibility as user_visibility, 
		IFNULL(u.profile_image, '') as profile_image, u.status as user_status, u.created_at as user_created_at, 
		u.updated_at as user_updated_at, u.updated_by as user_updated_by
		FROM users u
		WHERE u.status != 'delete'
		AND u.type != 'admin'
		AND u.id = ?
		ORDER BY u.id desc;
    `, user_id)
	if selectError != nil {
		return User{}, selectError
	}
	defer rows.Close()

	var user User

	for rows.Next() {
		// Scan the post and user data
		err := rows.Scan(
			&user.ID, &user.FirstName, &user.LastName, &user.NickName, &user.Email, &user.AboutMe, &user.Visibility,
			&user.ProfileImage, &user.Status, &user.CreatedAt,
			&user.UpdatedAt, &user.UpdatedBy,
		)
		if err != nil {
			return User{}, fmt.Errorf("error scanning row: %v", err)
		}

	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return User{}, fmt.Errorf("row iteration error: %v", err)
	}

	return user, nil
}

func UpdateStatusUser(db *dbTools.DBContainer, user_id int, status string, login_user_id int) error {
	updateQuery := `UPDATE users
					SET status = ?,
						updated_at = CURRENT_TIMESTAMP,
						updated_by = ?
					WHERE id = ?;`
	_, updateErr := db.Conn.Exec(updateQuery, status, login_user_id, user_id)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

func GetUserIDByUsername(db *dbTools.DBContainer, nick_name string) (int, error) {
	var userID int
	err := db.Conn.QueryRow("SELECT id FROM users WHERE nick_name = ?", nick_name).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, errors.New("user not found")
		}
		return -1, err // Other DB errors
	}

	return userID, nil
}
