package controller

import (
	"encoding/json"
	"fmt"
	"log"
)

func (cc *ChatController) listener() {
	for action := range cc.clientQueue {
		switch action.kind {
		case "join":
			cc.sendClientList(action.client, -1)

		case "online":
			cc.sendClientList(action.client, action.client.UserID)
			content := fmt.Sprintf(`{"action": "online", "id": "%d"}`, action.client.UserID)
			cc.queuePublicMessage(content, -1)

		case "offline":
			delete(cc.clients, action.client.UserID)
			content := fmt.Sprintf(`{"action": "offline", "id": "%d"}`, action.client.UserID)
			cc.queuePublicMessage(content, -1)
		}
	}
}

func (cc *ChatController) sendClientList(cl *client, receiver int) {
	type data struct {
		Action           string   `json:"action"`
		AllClients       []string `json:"allClients"`
		ClientIDs        []int    `json:"clientIDs"`
		OnlineClients    []int    `json:"onlineClients"`
		UnreadMsgClients []int    `json:"unreadMsgClients"`
	}

	d := data{Action: "start"}
	clientList, clientIDs, err := cc.db.SelectUserList(cl.UserID)
	if err != nil {
		log.Println("Error fetching UserList:", err)
		return
	}
	d.AllClients = *clientList
	d.ClientIDs = *clientIDs

	unreadList, err := cc.db.SelectUnreadMsgList(cl.UserID)
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
	cc.queuePublicMessage(string(jsonData), receiver)
}
