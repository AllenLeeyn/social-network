package messenger

import (
	"fmt"
	"log"
	"net/http"
	"social-network/pkg/dbTools"
	"time"

	"github.com/gorilla/websocket"
)

var db *dbTools.DBContainer
var upgrader = websocket.Upgrader{}

type message = dbTools.Message
type user = dbTools.User

type Messenger struct {
	clients     map[int]*client
	msgQueue    chan message
	clientQueue chan action
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

func Start(dbMain *dbTools.DBContainer) Messenger {
	db = dbMain
	m := Messenger{
		clients:     make(map[int]*client),
		msgQueue:    make(chan message, 100),
		clientQueue: make(chan action, 100),
	}
	go m.listener()
	go m.broadcaster()
	return m
}

func (m *Messenger) WebSocketUpgrade(w http.ResponseWriter, r *http.Request, sessionID string, u *user) {
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
	m.clients[u.ID] = cl

	if time.Since(u.RegDate) < time.Minute {
		m.clientQueue <- action{"join", cl}
	} else {
		m.clientQueue <- action{"online", cl}
	}
	go m.handleConnection(cl)
}

func (m *Messenger) queuePublicMessage(content string, tgt int) {
	m.msgQueue <- message{
		SenderID:   -1,
		ReceiverID: tgt,
		Content:    content,
	}
}

func (m *Messenger) sendMessage(msg message, tgt int) error {
	client, exists := m.clients[tgt]
	if exists {
		err := client.Conn.WriteMessage(websocket.TextMessage, []byte(msg.Content))
		if err != nil {
			return fmt.Errorf("error sending message to %v: %v", client, err)
		}
	}
	return nil
}

func (m *Messenger) CloseConn(s *dbTools.Session) error {
	u, err := db.SelectUserByField("id", s.UserID)
	if err != nil || u == nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	cl, exists := m.clients[u.ID]
	if !exists {
		return fmt.Errorf("user %v not found in clients map", u.ID)
	}
	cl.Conn.Close()
	return nil
}
