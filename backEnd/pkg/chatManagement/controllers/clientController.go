package controller

import (
	"encoding/json"
	"fmt"
	"log"

	chatModel "social-network/pkg/chatManagement/models"
	groupModel "social-network/pkg/groupManagement/models"
)

func (cc *ChatController) listener() {
	for action := range cc.clientQueue {
		switch action.kind {
		case "online":
			cc.sendClientList(action.client, action.client.UserUUID)
			status := fmt.Sprintf(`{"action": "online", "id": "%s", "name": "%s"}`, action.client.UserUUID, action.client.UserName)
			cc.queueStatusToFollowings(status, action.client.UserUUID)

		case "offline":
			delete(cc.clients, action.client.UserUUID)
			status := fmt.Sprintf(`{"action": "offline", "id": "%s"}`, action.client.UserUUID)
			cc.queueStatusToFollowings(status, action.client.UserUUID)
		}
	}
}

func (cc *ChatController) sendClientList(cl *client, receiverUUID string) {
	type data struct {
		Action              string                 `json:"action"`
		FollowingsName      []string               `json:"followingsName"`
		FollowingsUUID      []string               `json:"FollowingsUUID"`
		OnlineFollowings    []string               `json:"onlineFollowings"`
		UnreadMsgFollowings []string               `json:"unreadMsgFollowings"`
		GroupList           []groupModel.GroupView `json:"groupList"`
	}

	d := data{Action: "userList"}
	followingList, followingUUIDs, err := chatModel.SelectUserList(cl.UserID)
	if err != nil {
		log.Println("Error fetching UserList:", err)
		return
	}
	d.FollowingsName = *followingList
	d.FollowingsUUID = *followingUUIDs

	unreadList, err := chatModel.SelectUnreadMsgList(cl.UserID)
	if err != nil {
		log.Println("Error fetching UnreadMsgList:", err)
		return
	}
	d.UnreadMsgFollowings = *unreadList

	followings := make(map[string]bool)
	for _, uuid := range *followingUUIDs {
		followings[uuid] = true
	}
	for uuid := range cc.clients {
		if followings[uuid] {
			d.OnlineFollowings = append(d.OnlineFollowings, uuid)
		}
	}

	groups, err := groupModel.SelectGroups(cl.UserUUID, true)
	if err != nil {
		log.Println("Error fetching groupList:", err)
		return
	}
	d.GroupList = *groups

	jsonData, err := json.Marshal(d)
	if err != nil {
		log.Println("Error marshaling data to JSON:", err)
		return
	}
	cc.queuePublicMessage(string(jsonData), receiverUUID)
}
