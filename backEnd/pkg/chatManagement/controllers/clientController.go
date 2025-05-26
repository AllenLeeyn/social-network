package controller

import (
	"encoding/json"
	"fmt"
	"log"

	chatModel "social-network/pkg/chatManagement/models"
)

func (cc *ChatController) listener() {
	for action := range cc.clientQueue {
		switch action.kind {
		case "join":
			cc.sendClientList(action.client, "-1")

		case "online":
			cc.sendClientList(action.client, action.client.UserUUID)
			content := fmt.Sprintf(`{"action": "online", "id": "%s", "name": "%s"}`, action.client.UserUUID, action.client.UserName)
			cc.queuePublicMessage(content, "-1")

		case "offline":
			delete(cc.clients, action.client.UserUUID)
			content := fmt.Sprintf(`{"action": "offline", "id": "%s"}`, action.client.UserUUID)
			cc.queuePublicMessage(content, "-1")
		}
	}
}

func (cc *ChatController) sendClientList(cl *client, receiverUUID string) {
	type data struct {
		Action           string   `json:"action"`
		AllClients       []string `json:"allClients"`
		ClientUUIDs      []string `json:"clientIDs"`
		OnlineClients    []string `json:"onlineClients"`
		UnreadMsgClients []string `json:"unreadMsgClients"`
	}

	d := data{Action: "userList"}
	clientList, clientUUIDs, err := chatModel.SelectUserList(cl.UserID)
	if err != nil {
		log.Println("Error fetching UserList:", err)
		return
	}
	d.AllClients = *clientList
	d.ClientUUIDs = *clientUUIDs

	unreadList, err := chatModel.SelectUnreadMsgList(cl.UserID)
	if err != nil {
		log.Println("Error fetching UnreadMsgList:", err)
		return
	}
	d.UnreadMsgClients = *unreadList

	for userID := range cc.clients {
		d.OnlineClients = append(d.OnlineClients, userID)
	}

	jsonData, err := json.Marshal(d)
	if err != nil {
		log.Println("Error marshaling data to JSON:", err)
		return
	}
	cc.queuePublicMessage(string(jsonData), receiverUUID)
}
