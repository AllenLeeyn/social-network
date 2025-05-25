package models

import (
	"database/sql"
	"errors"
	"fmt"
	"social-network/pkg/utils"
	"time"
)

// User struct represents the user data model
type User struct {
	ID              int        `json:"id"`
	UUID            string     `json:"uuid"`
	TypeId          int        `json:"type_id"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	Gender          string     `json:"gender"`
	BirthDay        time.Time  `json:"birthday"`
	Email           string     `json:"email"`
	Password        string     `json:"password"`
	ConfirmPassword string     `json:"confirmPassword"`
	PasswordHash    string     `json:"pw_hash"`
	NickName        string     `json:"nick_name"`
	ProfileImage    string     `json:"profile_image"`
	AboutMe         string     `json:"about_me"`
	Visibility      string     `json:"visibility"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedBy       int        `json:"updated_by"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

type userView struct {
	UUID         string `json:"uuid"`
	NickName     string `json:"nick_name"`
	ProfileImage string `json:"profile_image"`
	Visibility   string `json:"visibility"`
}

type userProfile struct {
	UUID         string    `json:"uuid"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Gender       string    `json:"gender"`
	BirthDay     time.Time `json:"birthday"`
	Email        string    `json:"email"`
	NickName     string    `json:"nick_name"`
	ProfileImage string    `json:"profile_image"`
	AboutMe      string    `json:"about_me"`
	Visibility   string    `json:"visibility"`
}

func SelectUsers() (*[]userView, error) {
	qry := `SELECT uuid, nick_name, profile_image, visibility
			FROM users
			WHERE id != 0`

	rows, err := sqlDB.Query(qry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uViews []userView
	for rows.Next() {
		var uv userView
		err := rows.Scan(&uv.UUID, &uv.NickName, &uv.ProfileImage, &uv.Visibility)
		if err != nil {
			return nil, err
		}
		uViews = append(uViews, uv)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &uViews, nil
}

func SelectUser(tgtUUID string) (*userProfile, error) {
	qry := `SELECT
				uuid, first_name, last_name,
				gender, birthday,
				email, nick_name,
				profile_image, about_me, visibility
			FROM users
			WHERE uuid = ?
			LIMIT 1;`
	var uProfile userProfile
	err := sqlDB.QueryRow(qry, tgtUUID).Scan(
		&uProfile.UUID, &uProfile.FirstName, &uProfile.LastName,
		&uProfile.Gender, &uProfile.BirthDay,
		&uProfile.Email, &uProfile.NickName,
		&uProfile.ProfileImage, &uProfile.AboutMe, &uProfile.Visibility)
	if err != nil {
		return nil, err
	}
	return &uProfile, nil
}

func SelectUserIDByUUID(userUUID string) (int, error) {
	var userID int
	err := sqlDB.QueryRow(`SELECT id FROM users WHERE uuid = ?`, userUUID).Scan(&userID)
	return userID, err
}

func SelectUserByField(fieldName string, fieldValue interface{}) (*User, error) {
	if fieldName != "id" && fieldName != "nick_name" && fieldName != "email" && fieldName != "uuid"{
		return nil, fmt.Errorf("invalid field")
	}
	qry := `SELECT * FROM users WHERE ` + fieldName + ` = ?`

	var u User
	err := sqlDB.QueryRow(qry, fieldValue).Scan(
		&u.ID, &u.UUID, &u.TypeId,
		&u.FirstName, &u.LastName,
		&u.Gender, &u.BirthDay,
		&u.Email, &u.PasswordHash,
		&u.NickName, &u.ProfileImage, &u.AboutMe, &u.Visibility,
		&u.Status, &u.CreatedAt,
		&u.UpdatedBy, &u.UpdatedAt)
	if err != nil {
		return nil, checkErrNoRows(err)
	}
	return &u, nil
}

func checkUniqueUser(user *User) error {
	var existingEmail string
	var existingUsername string
	qry := `SELECT email, nick_name 
			FROM users 
			WHERE email = ? OR nick_name = ? LIMIT 1;`

	err := sqlDB.QueryRow(qry, user.Email, user.NickName).
		Scan(&existingEmail, &existingUsername)
	if err != nil {
		return checkErrNoRows(err)
	}

	if existingEmail == user.Email {
		return errors.New("email is already used")
	}
	if existingUsername == user.NickName {
		return errors.New("nick name is already used")
	}
	return nil
}

func InsertUser(user *User) (int, error) {
	uuid, err := utils.GenerateUuid()
	if err != nil {
		return -1, err
	}
	user.UUID = uuid

	if err = checkUniqueUser(user); err != nil {
		return -1, err
	}

	qry := `INSERT INTO users (
				uuid, type_id, 
				first_name, last_name, 
				birthday, gender, nick_name, 
				email, pw_hash,
				about_me, profile_image, visibility,
				updated_by, updated_at) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, CURRENT_TIMESTAMP);`
	result, insertErr := sqlDB.Exec(qry,
		user.UUID, 1,
		user.FirstName, user.LastName,
		user.BirthDay, user.Gender, user.NickName,
		user.Email, user.PasswordHash,
		user.AboutMe, user.ProfileImage, user.Visibility)

	if insertErr != nil {
		// Check if the error is a SQLite constraint violation (duplicate entry)
		if sqliteErr, ok := insertErr.(interface{ ErrorCode() int }); ok {
			if sqliteErr.ErrorCode() == 19 { // 19 = UNIQUE constraint failed (SQLite error code)
				return -1, errors.New("email or nick name already exists")
			}
		}
		return -1, insertErr // Other DB errors
	}

	userId, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(userId), nil
}

func UpdateUser(user *User) error {
	checkQuery := `SELECT id 
					FROM users 
					WHERE nick_name = ? AND id != ?`

	var existingID int
	err := sqlDB.QueryRow(checkQuery, user.NickName, user.ID).
		Scan(&existingID)
	if err == nil {
		return fmt.Errorf("nickname '%s' is already taken", user.NickName)
	}
	if err != sql.ErrNoRows {
		return err
	}

	updateQuery := `
		UPDATE users
		SET first_name = ?,	last_name = ?, nick_name =?,
			gender = ?, birthday = ?, about_me = ?,
			visibility = ?, profile_image = ?,
			status = ?, updated_at = CURRENT_TIMESTAMP, updated_by = ?
		WHERE id = ?;`

	_, err = sqlDB.Exec(updateQuery,
		user.FirstName, user.LastName, user.NickName,
		user.Gender, user.BirthDay, user.AboutMe,
		user.Visibility, user.ProfileImage,
		user.Status, user.UpdatedBy,
		user.ID,
	)

	return err
}

func IsPublic(userUUID string) bool {
	qry := `SELECT id 
			FROM users 
			WHERE uuid = ? AND visibility = 'public'`

	var existingID int
	err := sqlDB.QueryRow(qry, userUUID).Scan(&existingID)
	return err == nil
}
