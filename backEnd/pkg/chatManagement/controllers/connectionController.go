package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	errorControllers "social-network/pkg/errorManagement/controllers"
	middleware "social-network/pkg/middleware"
	userModel "social-network/pkg/userManagement/models"
)

func (cc *ChatController) WSHandler(w http.ResponseWriter, r *http.Request) {
	sessionId, userId, _, isOk := middleware.GetSessionCredentials(r.Context())
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	u, err := userModel.SelectUserByField("id", userId)
	if err != nil || u == nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	cc.WebSocketUpgrade(w, r, sessionId, u)
}

func (cc *ChatController) WebSocketUpgrade(w http.ResponseWriter, r *http.Request, sessionID string, u *user) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection: ", err)
		return
	}
	cl := &client{
		u.NickName,
		u.ID,
		u.UUID,
		conn,
	}
	cc.clients[u.UUID] = cl

	if time.Since(u.CreatedAt) < time.Minute {
		cc.clientQueue <- action{"join", cl}
	} else {
		cc.clientQueue <- action{"online", cl}
	}
	go cc.handleConnection(cl)
}

func (cc *ChatController) handleConnection(cl *client) {
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
			err = cc.processMessage(&msgData, cl)
		case "messageReq":
			err = cc.processMessageRequest(&msgData, cl)
		case "messageAck":
			err = cc.processMessageAcknowledgement(&msgData, cl)
		case "typing":
			cc.processTypingEvent(&msgData, cl)
		case "userListReq":
        	cc.sendClientList(cl, cl.UserUUID)
		}

		if err != nil {
			log.Println(err)
		}
	}
	cc.clientQueue <- action{"offline", cl}
	cl.Conn.Close()
}

func (cc *ChatController) CloseConn(userUUID string) error {
	cl, exists := cc.clients[userUUID]
	if !exists {
		return fmt.Errorf("user %v not found in clients map", userUUID)
	}
	cl.Conn.Close()
	return nil
}
