package controller

import (
	"log"

	"github.com/gorilla/websocket"
)

func (cc *ChatController) broadcaster() {
	for {
		msg := <-cc.msgQueue

		if msg.ReceiverID == -1 {
			for _, client := range cc.clients {
				err := client.Conn.WriteMessage(websocket.TextMessage, []byte(msg.Content))
				if err != nil {
					log.Printf("Error sending message to %v: %v", msg.ReceiverID, err)
				}
			}
		} else {
			err := cc.sendMessage(msg, msg.ReceiverID)
			if err != nil {
				log.Printf("Error sending message to %v: %v", msg.ReceiverID, err)
			}
			if msg.SenderID != -1 && msg.SenderID != msg.ReceiverID {
				err = cc.sendMessage(msg, msg.SenderID)
				if err != nil {
					log.Printf("Error sending message to %v: %v", msg.SenderID, err)
				}
			}
		}
	}
}
