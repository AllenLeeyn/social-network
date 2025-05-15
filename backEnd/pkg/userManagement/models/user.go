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
	ID              int            `json:"id"`
	UUID            string         `json:"uuid"`
	TypeId          int            `json:"type_id"`
	FirstName       string         `json:"first_name"`
	LastName        string         `json:"last_name"`
	Gender          string         `json:"gender"`
	BirthDay        time.Time      `json:"birthday"`
	Email           string         `json:"email"`
	Password        string         `json:"password"`
	ConfirmPassword string         `json:"confirmPassword"`
	PasswordHash    []byte         `json:"pw_hash"`
	NickName        string         `json:"nick_name"`
	ProfileImage    sql.NullString `json:"profile_image"`
	AboutMe         string         `json:"about_me"`
	Visibility      string         `json:"visibility"`
	Status          string         `json:"status"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedBy       int            `json:"updated_by"`
	UpdatedAt       *time.Time     `json:"updated_at"`
}

func SelectUserByField(fieldName string, fieldValue interface{}) (*User, error) {
	if fieldName != "id" && fieldName != "nick_name" && fieldName != "email" {
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

	err := sqlDB.QueryRow(qry, user.Email, user.NickName).Scan(&existingEmail, &existingUsername)
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
	updateQuery := `
		UPDATE users
		SET first_name = ?,	last_name = ?, nick_name =?,
			gender = ?, birthday = ?, about_me = ?,
			visibility = ?, profile_image = ?,
			status = ?, updated_at = CURRENT_TIMESTAMP, updated_by = ?
		WHERE id = ?;`

	_, err := sqlDB.Exec(updateQuery,
		user.FirstName, user.LastName, user.NickName,
		user.Gender, user.BirthDay, user.AboutMe,
		user.Visibility, user.ProfileImage,
		user.Status, user.UpdatedBy,
		user.ID,
	)

	return err
}
