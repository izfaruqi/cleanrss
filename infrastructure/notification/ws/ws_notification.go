package ws

import (
	"cleanrss/domain"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type wsClient struct{}
type wsNotificationService struct {
	mu          sync.Mutex
	subscribers []func(notification domain.Notification)

	wsClients          map[*websocket.Conn]wsClient
	wsRegisterClient   chan *websocket.Conn
	wsNotifications    chan domain.Notification
	wsUnregisterClient chan *websocket.Conn
}

func NewWSNotificationService() (domain.NotificationService, http.Handler) {
	ns := &wsNotificationService{wsClients: make(map[*websocket.Conn]wsClient), wsRegisterClient: make(chan *websocket.Conn), wsNotifications: make(chan domain.Notification), wsUnregisterClient: make(chan *websocket.Conn)}
	router := chi.NewRouter()
	go func() {
		err := ns.initWS()
		if err != nil {
			log.Println(err)
		}
	}()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024, CheckOrigin: func(r *http.Request) bool { // TODO: Only allow connections from CleanRSS.
			return true
		}}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}
		ns.handleWS(conn)
	})
	return ns, router
}

func (w *wsNotificationService) initWS() error {
	for {
		select {
		case wsConn := <-w.wsRegisterClient:
			w.wsClients[wsConn] = wsClient{}
			log.Println("Connection registered.")
		case notification := <-w.wsNotifications:
			notificationJson, _ := json.Marshal(notification)
			for wsConn := range w.wsClients {
				if err := wsConn.WriteMessage(websocket.TextMessage, notificationJson); err != nil {
					log.Println("Write error. Closing connection.")
					w.wsUnregisterClient <- wsConn
					err := wsConn.WriteMessage(websocket.CloseMessage, []byte{})
					if err != nil {
						return err
					}
					err = wsConn.Close()
					if err != nil {
						return err
					}
				}
			}
		case wsConn := <-w.wsUnregisterClient:
			delete(w.wsClients, wsConn)
			log.Println("Connection unregistered.")
		}
	}
}

func (w *wsNotificationService) handleWS(c *websocket.Conn) {
	defer func() {
		w.wsUnregisterClient <- c
		err := c.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}()
	w.wsRegisterClient <- c

	var (
		mt  int
		msg []byte
		err error
	)
	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			log.Println("read:", err)
			break
		}
		w.forward(domain.Notification{Code: "WS", Payload: string(msg)})
		if err = c.WriteMessage(mt, msg); err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func (w *wsNotificationService) forward(notification domain.Notification) {
	for _, handler := range w.subscribers {
		go handler(notification)
	}
}

func (w *wsNotificationService) Publish(notification domain.Notification) {
	w.wsNotifications <- notification
}

func (w *wsNotificationService) Subscribe(handler func(notification domain.Notification)) {
	w.mu.Lock()
	w.subscribers = append(w.subscribers, handler)
	w.mu.Unlock()
}
