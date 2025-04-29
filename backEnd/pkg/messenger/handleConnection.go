package messenger

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"strings"
	"time"
)

func (m *Messenger) handleConnection(cl *client) {
	for {
		_, msg, err := cl.Conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		msgData := message{}
		err = json.Unmarshal(msg, &msgData)
		if err != nil {
			log.Printf("Invalid message format")
			continue
		}

		switch msgData.Action {
		case "message":
			err = m.processMessage(&msgData, cl)
		case "messageReq":
			err = m.processMessageRequest(&msgData, cl)
		case "messageAck":
			err = m.processMessageAcknowledgement(&msgData, cl)
		case "typing":
			m.processTypingEvent(&msgData, cl)
		}

		if err != nil {
			log.Println(err)
		}
	}
	m.clientQueue <- action{"offline", cl}
	cl.Conn.Close()
}

func (m *Messenger) processMessage(msgData *message, cl *client) error {
	isValidMsg, sanitizeMsg := checkMessage(msgData.Content)
	if !isValidMsg {
		return fmt.Errorf("invalid message")
	}

	receiver, err := db.SelectUserByField("id", msgData.ReceiverID)
	if err != nil || receiver == nil || receiver.ID == 0 {
		return fmt.Errorf("receiver not found: %v", err)
	}

	msgData.SenderID = cl.UserID
	msgData.Content = sanitizeMsg
	msgData.CreatedAt = time.Now()

	err = db.InsertMessage(msgData)
	if err != nil {
		return fmt.Errorf("receiver inserting message: %v", err)
	}

	content, err := json.Marshal(msgData)
	if err != nil {
		log.Printf("Error generating JSON: %v", err)
	}
	msgData.Content = string(content)

	m.queuePublicMessage(`{"action":"messageSendOK"}`, cl.UserID)
	m.msgQueue <- *msgData
	return nil
}

func (m *Messenger) processMessageRequest(msgData *message, cl *client) error {
	msgIdStr := msgData.Content
	messages, err := db.SelectMessages(cl.UserID, msgData.ReceiverID, msgIdStr)
	if err != nil {
		return fmt.Errorf("error getting messages: %v", err)
	}

	content := struct {
		Action  string    `json:"action"`
		Content []message `json:"content"`
	}{Action: "messageHistory", Content: *messages}

	contentJSON, err := json.Marshal(content)
	if err != nil {
		log.Printf("Error generating JSON: %v", err)
	}
	m.queuePublicMessage(string(contentJSON), cl.UserID)
	return nil
}

func (m *Messenger) processMessageAcknowledgement(msgData *message, cl *client) error {
	messages, err := db.SelectUnreadMessages(msgData.ReceiverID, cl.UserID)
	if err != nil {
		return fmt.Errorf("error getting messages: %v", err)
	}

	for _, msg := range *messages {
		msg.ReadAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		err := db.UpdateMessage(&msg)
		if err != nil {
			return fmt.Errorf("error acknowledging messages: %v", err)
		}
	}
	return nil
}

func checkMessage(message string) (bool, string) {
	message = strings.TrimSpace(message)
	if len(message) == 0 {
		return false, "Too short"
	} else if len(message) > 1000 {
		return false, "Too long"
	}
	return true, html.EscapeString(message)
}

func (m *Messenger) processTypingEvent(msgData *message, cl *client) {
	content := fmt.Sprintf(`{"action": "typing", "receiverID": %d, "senderID": %d}`, msgData.ReceiverID, cl.UserID)
	m.queuePublicMessage(content, msgData.ReceiverID)
}
