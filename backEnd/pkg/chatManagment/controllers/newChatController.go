package controller

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	chatModel "social-network/pkg/chatManagment/models"
	"social-network/pkg/dbTools"
	userModel "social-network/pkg/userManagement/models"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type message = chatModel.Message
type user = userModel.User

type ChatController struct {
	clients     map[int]*client
	msgQueue    chan message
	clientQueue chan action
	db          *chatModel.ChatModel
	um          *userModel.UserModel
}

type client struct {
	UserName  string
	UserID    int
	SessionID string
	Conn      *websocket.Conn
}

type action struct {
	kind   string
	client *client
}

func NewChatController(dbMain *sql.DB) *ChatController {
	cc := &ChatController{
		clients:     make(map[int]*client),
		msgQueue:    make(chan message, 100),
		clientQueue: make(chan action, 100),
		db:          chatModel.NewChatModel(dbMain),
		um:          userModel.NewUserModel(dbMain),
	}
	go cc.listener()
	go cc.broadcaster()
	return cc
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
		sessionID,
		conn,
	}
	cc.clients[u.ID] = cl

	if time.Since(u.CreatedAt) < time.Minute {
		cc.clientQueue <- action{"join", cl}
	} else {
		cc.clientQueue <- action{"online", cl}
	}
	go cc.handleConnection(cl)
}

func (cc *ChatController) queuePublicMessage(content string, tgt int) {
	cc.msgQueue <- message{
		SenderID:   -1,
		ReceiverID: tgt,
		Content:    content,
	}
}

func (cc *ChatController) sendMessage(msg message, tgt int) error {
	client, exists := cc.clients[tgt]
	if exists {
		err := client.Conn.WriteMessage(websocket.TextMessage, []byte(msg.Content))
		if err != nil {
			return fmt.Errorf("error sending message to %v: %v", client, err)
		}
	}
	return nil
}

func (cc *ChatController) CloseConn(s *dbTools.Session) error {
	u, err := cc.um.SelectUserByField("id", s.UserID)
	if err != nil || u == nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	cl, exists := cc.clients[u.ID]
	if !exists {
		return fmt.Errorf("user %v not found in clients map", u.ID)
	}
	cl.Conn.Close()
	return nil
}

func (cc *ChatController) SelectUserByField(fieldName string, fieldValue interface{}) (*user, error) {
	return cc.um.SelectUserByField(fieldName, fieldValue)
}
