package ws

import (
	"encoding/json"
	"log"

	"github.com/gofiber/websocket/v2"
)

type client struct{}
type Notification struct {
	Code    string `json:"code"`
	Payload string `json:"payload"`
}

var (
	wsClients          = make(map[*websocket.Conn]client)
	wsRegisterClient   = make(chan *websocket.Conn)
	WSNotifications    = make(chan Notification)
	wsUnregisterClient = make(chan *websocket.Conn)
)

func WSInit() {
	for {
		select {
		case wsConn := <-wsRegisterClient:
			wsClients[wsConn] = client{}
			log.Println("Connection registered.")
		case notification := <-WSNotifications:
			log.Println("Notification: " + notification.Code + " " + notification.Payload)
			notificationJson, _ := json.Marshal(notification)
			for wsConn := range wsClients {
				if err := wsConn.WriteMessage(websocket.TextMessage, notificationJson); err != nil {
					log.Println("Write error. Closing connection.")
					wsUnregisterClient <- wsConn
					wsConn.WriteMessage(websocket.CloseMessage, []byte{})
					wsConn.Close()
				}
			}
		case wsConn := <-wsUnregisterClient:
			delete(wsClients, wsConn)
			log.Println("Connection unregistered.")
		}
	}
}

func WSController(c *websocket.Conn) {

	defer func() {
		wsUnregisterClient <- c
		c.Close()
	}()
	wsRegisterClient <- c

	var (
		mt  int
		msg []byte
		err error
	)
	for {
		log.Println("WS client connected.")
		if mt, msg, err = c.ReadMessage(); err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", msg)
		if err = c.WriteMessage(mt, msg); err != nil {
			log.Println("write:", err)
			break
		}
	}
}
