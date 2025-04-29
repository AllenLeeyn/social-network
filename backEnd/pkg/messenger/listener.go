package messenger

import (
	"encoding/json"
	"fmt"
	"log"
)

func (m *Messenger) listener() {
	for action := range m.clientQueue {
		switch action.kind {
		case "join":
			m.sendClientList(action.client, -1)

		case "online":
			m.sendClientList(action.client, action.client.UserID)
			content := fmt.Sprintf(`{"action": "online", "id": "%d"}`, action.client.UserID)
			m.queuePublicMessage(content, -1)

		case "offline":
			delete(m.clients, action.client.UserID)
			content := fmt.Sprintf(`{"action": "offline", "id": "%d"}`, action.client.UserID)
			m.queuePublicMessage(content, -1)
		}
	}
}

func (m *Messenger) sendClientList(cl *client, receiver int) {
	type data struct {
		Action           string   `json:"action"`
		AllClients       []string `json:"allClients"`
		ClientIDs        []int    `json:"clientIDs"`
		OnlineClients    []int    `json:"onlineClients"`
		UnreadMsgClients []int    `json:"unreadMsgClients"`
	}

	d := data{Action: "start"}
	clientList, clientIDs, err := db.SelectUserList(cl.UserID)
	if err != nil {
		log.Println("Error fetching UserList:", err)
		return
	}
	d.AllClients = *clientList
	d.ClientIDs = *clientIDs

	unreadList, err := db.SelectUnreadMsgList(cl.UserID)
	if err != nil {
		log.Println("Error fetching UnreadMsgList:", err)
		return
	}
	d.UnreadMsgClients = *unreadList

	for userID := range m.clients {
		d.OnlineClients = append(d.OnlineClients, userID)
	}

	jsonData, err := json.Marshal(d)
	if err != nil {
		log.Println("Error marshaling data to JSON:", err)
		return
	}
	m.queuePublicMessage(string(jsonData), receiver)
}
