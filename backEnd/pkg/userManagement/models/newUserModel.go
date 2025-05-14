package models

import "database/sql"

type UserModel struct {
	DB *sql.DB
}

func NewUserModel(dbMain *sql.DB) *UserModel {
	return &UserModel{DB: dbMain}
}

// checkErrNoRows() checks if no result from sql query.
func checkErrNoRows(err error) error {
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}
