package controller

import (
	"log"
	chatModel "social-network/pkg/chatManagement/models"
	userModel "social-network/pkg/userManagement/models"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type message = chatModel.Message
type user = userModel.User

type ChatController struct {
	clients     map[string]*client
	msgQueue    chan message
	clientQueue chan action
}

type client struct {
	UserName string
	UserID   int
	UserUUID string
	Conn     *websocket.Conn
}

type action struct {
	kind   string
	client *client
}

func Initialize() *ChatController {
	log.Println("\033[35mInitlise chat controller\033[0m")
	cc := &ChatController{
		clients:     make(map[string]*client),
		msgQueue:    make(chan message, 100),
		clientQueue: make(chan action, 100),
	}
	go cc.listener()
	go cc.broadcaster()
	return cc
}
