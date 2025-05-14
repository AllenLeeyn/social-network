package model

import "database/sql"

type ChatModel struct {
	DB *sql.DB
}

func NewChatModel(dbMain *sql.DB) *ChatModel {
	return &ChatModel{DB: dbMain}
}

// checkErrNoRows() checks if no result from sql query.
func checkErrNoRows(err error) error {
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}
