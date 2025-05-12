package models

import (
	"database/sql"
	"forum/db"
	"log"
	"social-network/pkg/utils"
	"time"
)

// Chat represents the "chats" table
type Chat struct {
	ID        int       `json:"id"`
	UUID      string    `json:"uuid"`
	Type      string    `json:"type"` // "private" or "group"
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy int       `json:"created_by"`
}

// ChatMember represents the "chat_members" table
type ChatMember struct {
	ID        int       `json:"id"`
	ChatID    int       `json:"chat_id"`
	UserID    int       `json:"user_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// Message represents the "messages" table
type Message struct {
	ID                int        `json:"id"`
	ChatID            int        `json:"chat_id"`
	Content           string     `json:"content"`
	Status            string     `json:"status"`
	CreatedAt         time.Time  `json:"created_at"`
	CreatedBy         int        `json:"created_by"`
	UpdatedAt         *time.Time `json:"updated_at"`
	UpdatedBy         *int       `json:"updated_by"`
	CreatedByUsername string     `json:"created_by_username"`
}

// MessageFile represents the "message_files" table
type MessageFile struct {
	ID               int        `json:"id"`
	ChatID           int        `json:"chat_id"`
	MessageID        int        `json:"message_id"`
	FileUploadedName string     `json:"file_uploaded_name"`
	FileRealName     string     `json:"file_real_name"`
	Status           string     `json:"status"`
	CreatedAt        time.Time  `json:"created_at"`
	CreatedBy        int        `json:"created_by"`
	UpdatedAt        *time.Time `json:"updated_at"`
	UpdatedBy        *int       `json:"updated_by"`
}

func CheckChatExists(user1ID, user2ID int) (int, error) {
	db := db.OpenDBConnection()
	defer db.Close()

	var chatID int

	query := `
        SELECT c.id
        FROM chats c
        JOIN chat_members cm1 ON c.id = cm1.chat_id
        JOIN chat_members cm2 ON c.id = cm2.chat_id
        WHERE c.type = 'private' AND cm1.user_id = ? AND cm2.user_id = ?
    `
	err := db.QueryRow(query, user1ID, user2ID).Scan(&chatID)
	if err == sql.ErrNoRows {
		return 0, nil // No chat exists
	} else if err != nil {
		// Unexpected error
		return 0, err
	}

	return chatID, nil
}

func InsertChat(chat *Chat, user1ID, user2ID int, uploadedFiles map[string]string) (int, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Start a transaction for atomicity
	tx, err := db.Begin()
	if err != nil {
		return -1, err
	}

	chat.UUID, err = utils.GenerateUuid()
	if err != nil {
		tx.Rollback() // Rollback if UUID generation fails
		return -1, err
	}

	insertQuery := `INSERT INTO chats (uuid, type, created_by) VALUES (?, ?, ?);`
	result, insertErr := tx.Exec(insertQuery, chat.UUID, "private", user1ID) // Assuming user1ID initiated chat
	if insertErr != nil {
		return -1, insertErr
	}

	// Retrieve the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return -1, err
	}

	insertChatMember1Err := InsertChatMember(int(lastInsertID), user1ID, tx)
	if insertChatMember1Err != nil {
		tx.Rollback()
		return -1, insertChatMember1Err
	}
	insertChatMember2Err := InsertChatMember(int(lastInsertID), user2ID, tx)
	if insertChatMember2Err != nil {
		tx.Rollback()
		return -1, insertChatMember2Err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return -1, err
	}

	return int(lastInsertID), nil
}

func InsertChatMember(chatID, userID int, tx *sql.Tx) error {
	insertChatMemberQuery := `INSERT INTO chat_members (chat_id, user_id) VALUES (?, ?);`
	_, err := tx.Exec(insertChatMemberQuery, chatID, userID)
	if err != nil {
		return err
	}
	return nil
}

func InsertMsg(msg *Message, uploadedFiles map[string]string) (int, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Start a transaction for atomicity
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting InsertMsg transaction: %v", err)
		return -1, err
	}

	msg.Content = utils.SanitizeInput(msg.Content)

	insertMsgQuery := `INSERT INTO messages (chat_id, content, status, created_at, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?);`
	result, insertMsgQueryErr := tx.Exec(insertMsgQuery, msg.ChatID, msg.Content, msg.Status, msg.CreatedAt, msg.CreatedBy, msg.UpdatedBy)
	if insertMsgQueryErr != nil {
		tx.Rollback()
		log.Printf("Error inserting message: %v", insertMsgQueryErr)
		return -1, insertMsgQueryErr
	}

	// Retrieve the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Printf("Error retrieving last insert ID: %v", err)
		return -1, err
	}

	// Update the chats table's updated_at and updated_by fields
	updateChatQuery := `UPDATE chats SET updated_at = ?, updated_by = ? WHERE id = ?;`
	_, updateChatErr := tx.Exec(updateChatQuery, msg.CreatedAt, msg.CreatedBy, msg.ChatID)
	if updateChatErr != nil {
		tx.Rollback()
		log.Printf("Error updating chat: %v", updateChatErr)
		return -1, updateChatErr
	}

	// Check if there are uploaded files before calling InsertMsgFiles
	if len(uploadedFiles) > 0 {
		insertMsgFilesErr := InsertMsgFiles(msg.ChatID, int(lastInsertID), uploadedFiles, msg.CreatedBy, tx)
		if insertMsgFilesErr != nil {
			tx.Rollback()
			log.Printf("Error inserting message files: %v", insertMsgFilesErr)
			return -1, insertMsgFilesErr
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		log.Printf("Error committing transaction: %v", err)
		return -1, err
	}

	return int(lastInsertID), nil
}

func InsertMsgFiles(chatID int, msgID int, uploadedFiles map[string]string, user_id int, tx *sql.Tx) error {
	if len(uploadedFiles) > 0 {
		query := `INSERT INTO message_files (chat_id, message_id, file_real_name, file_uploaded_name, created_by) VALUES `
		values := make([]any, 0, len(uploadedFiles)*3)

		for i := range len(uploadedFiles) {
			if i > 0 {
				query += ", "
			}
			query += "(?, ?, ?, ?, ?)"
		}
		for key, value := range uploadedFiles {
			values = append(values, chatID, msgID, key, value, user_id)
		}
		query += ";"

		// Execute the bulk insert query
		_, err := tx.Exec(query, values...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return nil
}

func ReadAllMsgs(chatID, userID int) ([]Message, error) {
	db := db.OpenDBConnection()
	defer db.Close()

	var messages []Message

	query := `
		SELECT 
			m.id, m.chat_id, m.content, m.status, m.created_at, m.created_by, 
			m.updated_at, m.updated_by, u.username AS created_by_username
		FROM messages m
		JOIN chat_members cm ON m.chat_id = cm.chat_id
		JOIN users u ON m.created_by = u.id
		WHERE m.chat_id = ? AND cm.user_id = ?;
	`
	rows, err := db.Query(query, chatID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var message Message
		if err := rows.Scan(
			&message.ID,
			&message.ChatID,
			&message.Content,
			&message.Status,
			&message.CreatedAt,
			&message.CreatedBy,
			&message.UpdatedAt,
			&message.UpdatedBy,
			&message.CreatedByUsername,
		); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}
