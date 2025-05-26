package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	chatModel "social-network/pkg/chatManagement/models"
	userModel "social-network/pkg/userManagement/models"
	"social-network/pkg/utils"

	"github.com/gorilla/websocket"
)

func (cc *ChatController) broadcaster() {
	for {
		msg := <-cc.msgQueue

		if msg.ReceiverUUID == "-1" {
			for _, client := range cc.clients {
				err := client.Conn.WriteMessage(websocket.TextMessage, []byte(msg.Content))
				if err != nil {
					log.Printf("Error sending message to %v: %v", msg.ReceiverID, err)
				}
			}
		} else {
			log.Printf("broadcaster: Attempting to send to ReceiverUUID: %s", msg.ReceiverUUID)
			err := cc.sendMessage(msg, msg.ReceiverUUID)
			if err != nil {
				log.Printf("Error sending message to %v: %v", msg.ReceiverID, err)
			}
			if msg.SenderID != -1 && msg.SenderID != msg.ReceiverID {
				err = cc.sendMessage(msg, msg.SenderUUID)
				if err != nil {
					log.Printf("Error sending message to %v: %v", msg.SenderID, err)
				}
			}
		}
	}
}

func (cc *ChatController) queuePublicMessage(content string, tgt string) {
	cc.msgQueue <- message{
		SenderID:     -1,
		ReceiverUUID: tgt,
		Content:      content,
	}
}

func (cc *ChatController) sendMessage(msg message, tgt string) error {
	client, exists := cc.clients[tgt]
	if exists {
		err := client.Conn.WriteMessage(websocket.TextMessage, []byte(msg.Content))
		if err != nil {
			return fmt.Errorf("error sending message to %v: %v", client, err)
		}
	}
	return nil
}

func (cc *ChatController) processMessage(msgData *message, cl *client) error {
	sanitizeMsg, isValidMsg := checkMessage(msgData.Content)
	if !isValidMsg {
		return fmt.Errorf("invalid message")
	}

	receiver, err := userModel.SelectUserByField("uuid", msgData.ReceiverUUID)
	if err != nil || receiver == nil || receiver.ID == 0 {
		return fmt.Errorf("receiver not found: %v", err)
	}

	msgData.SenderID = cl.UserID
	msgData.ReceiverID = receiver.ID
	msgData.Content = sanitizeMsg
	msgData.CreatedAt = time.Now()

	err = chatModel.InsertMessage(msgData)
	if err != nil {
		return fmt.Errorf("receiver inserting message: %v", err)
	}

	content, err := json.Marshal(msgData)
	if err != nil {
		log.Printf("Error generating JSON: %v", err)
	}
	msgData.Content = string(content)

	cc.queuePublicMessage(`{"action":"messageSendOK"}`, cl.UserUUID)
	cc.msgQueue <- *msgData
	return nil
}

func (cc *ChatController) processMessageRequest(msgData *message, cl *client) error {
	msgIdStr := msgData.Content
	receiver, err := userModel.SelectUserByField("uuid", msgData.ReceiverUUID)
	if err != nil || receiver == nil {
		return fmt.Errorf("receiver not found: %v", err)
	}
	messages, err := chatModel.SelectMessages(cl.UserID, receiver.ID, msgIdStr)
	if err != nil {
		return fmt.Errorf("error getting messages: %v", err)
	}
	// log.Printf("processMessageRequest: Fetched messages: %+v", messages)

	content := struct {
		Action  string    `json:"action"`
		Content []message `json:"content"`
	}{Action: "messageHistory", Content: *messages}

	contentJSON, err := json.Marshal(content)
	if err != nil {
		log.Printf("Error generating JSON: %v", err)
	}
	// log.Printf("processMessageRequest: Sending messageHistory JSON: %s", string(contentJSON))

	cc.queuePublicMessage(string(contentJSON), cl.UserUUID)
	return nil
}

func (cc *ChatController) processMessageAcknowledgement(msgData *message, cl *client) error {
	receiver, err := userModel.SelectUserByField("uuid", msgData.ReceiverUUID)
	if err != nil || receiver == nil {
		return fmt.Errorf("receiver not found: %v", err)
	}
	messages, err := chatModel.SelectUnreadMessages(receiver.ID, cl.UserID)
	if err != nil {
		return fmt.Errorf("error getting messages: %v", err)
	}

	for _, msg := range *messages {
		msg.ReadAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		err := chatModel.UpdateMessage(&msg)
		if err != nil {
			return fmt.Errorf("error acknowledging messages: %v", err)
		}
	}

	cc.queuePublicMessage(
		fmt.Sprintf(`{"action":"messageAck", "senderUUID":"%s"}`, cl.UserUUID),
		msgData.ReceiverUUID,
	)
	return nil
}

func checkMessage(message string) (string, bool) {
	return utils.IsValidContent(message, 1, 1000)
}

func (cc *ChatController) processTypingEvent(msgData *message, cl *client) {
	content := fmt.Sprintf(
		`{"action": "typing", "receiverUUID": "%s", "senderUUID": "%s"}`,
		msgData.ReceiverUUID,
		cl.UserUUID,
	)
	cc.queuePublicMessage(content, msgData.ReceiverUUID)
}

