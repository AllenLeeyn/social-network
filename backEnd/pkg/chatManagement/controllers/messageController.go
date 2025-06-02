package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	chatModel "social-network/pkg/chatManagement/models"
	followingModel "social-network/pkg/followingManagement/models"
	groupModel "social-network/pkg/groupManagement/models"
	"social-network/pkg/utils"

	"github.com/gorilla/websocket"
)

const publicGroupUUID = "00000000-0000-0000-0000-000000000000"

func (cc *ChatController) broadcaster() {
	for {
		msg := <-cc.msgQueue

		switch {
		case msg.ReceiverUUID == "-1":
			cc.broadcastToAll(msg)
		case msg.GroupUUID != publicGroupUUID:
			cc.broadcastToGroup(msg)
		default:
			cc.broadcastPrivate(msg)
		}
	}
}

func (cc *ChatController) broadcastToAll(msg message) {
	for _, client := range cc.clients {
		if err := client.Conn.WriteMessage(websocket.TextMessage, []byte(msg.Content)); err != nil {
			log.Printf("Error sending message to %v: %v", client.UserUUID, err)
		}
	}
}

func (cc *ChatController) broadcastPrivate(msg message) {
	log.Printf("broadcaster: Attempting to send to ReceiverUUID: %s", msg.ReceiverUUID)
	err := cc.sendMessage(msg, msg.ReceiverUUID)
	if err != nil {
		log.Printf("Error sending message to %v: %v", msg.ReceiverID, err)
	}
	if msg.SenderID != -1 && msg.SenderUUID != msg.ReceiverUUID {
		err = cc.sendMessage(msg, msg.SenderUUID)
		if err != nil {
			log.Printf("Error sending message to %v: %v", msg.SenderID, err)
		}
	}
}

func (cc *ChatController) broadcastToGroup(msg message) {
	members, err := groupModel.SelectGroupMembers(msg.GroupUUID, "accepted")
	log.Printf("broadcaster: Sending group message to %d members", len(*members))
	if err != nil {
		log.Printf("Failed to get group members: %v", err)
		return
	}
	for _, m := range *members {
		_ = cc.sendMessage(msg, m.FollowerUUID)
	}
}

func (cc *ChatController) queuePublicMessage(content string, tgt string) {
	cc.msgQueue <- message{
		SenderID:     -1,
		ReceiverUUID: tgt,
		GroupUUID:    publicGroupUUID,
		Content:      content,
	}
}

func (cc *ChatController) queueStatusToFollowings(status string, userUUID string) {
	followings, err := followingModel.SelectFollowings(userUUID, "accepted")
	if err != nil {
		log.Printf("Failed to get followings: %v", err)
		return
	}

	for _, f := range *followings {
		tgt := f.FollowerUUID
		if tgt == userUUID {
			tgt = f.LeaderUUID
		}
		cc.msgQueue <- message{
			SenderID:     -1,
			ReceiverUUID: tgt,
			Content:      status,
		}
	}
}

func (cc *ChatController) sendMessage(msg message, tgt string) error {
	client, exists := cc.clients[tgt]
	if exists {
		log.Println("Send to : ", tgt)
		err := client.Conn.WriteMessage(websocket.TextMessage, []byte(msg.Content))
		if err != nil {
			return fmt.Errorf("error sending message to %v: %v", client, err)
		}
	} else {
		log.Println("unable to send to : ", tgt)
	}
	return nil
}

func (cc *ChatController) processRequestPermission(msgData *message, cl *client) error {
	if msgData.GroupUUID != publicGroupUUID {
		if !groupModel.IsGroupMember(msgData.GroupUUID, cl.UserID) {
			return fmt.Errorf("user not group member")
		}
	} else {
		if !followingModel.IsFollower(cl.UserID, msgData.ReceiverUUID) &&
			!followingModel.IsLeader(cl.UserID, msgData.ReceiverUUID) {
			return fmt.Errorf("user not part of following")
		}
	}
	return nil
}

func (cc *ChatController) processMessage(msgData *message, cl *client) error {
	if err := cc.processRequestPermission(msgData, cl); err != nil {
		return err
	}

	sanitizeMsg, isValidMsg := checkMessage(msgData.Content)
	if !isValidMsg {
		return fmt.Errorf("invalid message")
	}
	msgData.Content = sanitizeMsg
	msgData.SenderUUID = cl.UserUUID
	if msgData.GroupUUID == "" {
		msgData.GroupUUID = publicGroupUUID
	}

	msgID, err := chatModel.InsertMessage(msgData)
	if err != nil {
		return fmt.Errorf("receiver inserting message: %v", err)
	}

	newMsg, err := chatModel.SelectMessage(msgID)
	if err != nil {
		log.Printf("Error getting new msg: %v", err)
	}

	content, err := json.Marshal(newMsg)
	if err != nil {
		log.Printf("Error generating JSON: %v", err)
	}
	msgData.Content = string(content)

	cc.queuePublicMessage(`{"action":"messageSendOK"}`, cl.UserUUID)
	cc.msgQueue <- *msgData
	return nil
}

func (cc *ChatController) processMessageRequest(msgData *message, cl *client) error {
	log.Printf("procssing message request")
	if err := cc.processRequestPermission(msgData, cl); err != nil {
		return err
	}

	msgIdStr := msgData.Content
	operation := chatModel.SelectPrivateMessages
	params := []string{cl.UserUUID, msgData.ReceiverUUID}

	if msgData.GroupUUID != publicGroupUUID {
		operation = chatModel.SelectGroupMessages
		params = []string{"", msgData.GroupUUID}
	}

	messages, err := operation(params[0], params[1], msgIdStr)
	if err != nil {
		return fmt.Errorf("error getting messages: %v", err)
	}

	content := struct {
		Action  string        `json:"action"`
		Content []messageView `json:"content"`
	}{Action: "messageHistory", Content: *messages}

	contentJSON, err := json.Marshal(content)
	if err != nil {
		log.Printf("Error generating JSON: %v", err)
	}

	cc.queuePublicMessage(string(contentJSON), cl.UserUUID)
	return nil
}

func (cc *ChatController) processMessageAcknowledgement(msgData *message, cl *client) error {
	messages, err := chatModel.SelectUnreadMessages(msgData.ReceiverUUID, cl.UserUUID)
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
