package service

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/yosev/debugo"
)

type WebsocketService struct {
}

type Subject = string

const (
	TASK_CREATED Subject = "task:created"
	TASK_DELETED Subject = "task:deleted"
	TASK_UPDATED Subject = "task:updated"
)

type message struct {
	Subject string `json:"subject"`
	Payload any    `json:"payload"`
}

var (
	debug = debugo.New("websocket:service")
	mutex = sync.RWMutex{}
)

var conns = make(map[string]*websocket.Conn)

func (s *WebsocketService) AddConnection(uuid string, conn *websocket.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	conns[uuid] = conn
}

func (s *WebsocketService) RemoveConnection(uuid string, conn *websocket.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(conns, uuid)
}

func (s *WebsocketService) Broadcast(subject Subject, msg any) error {
	mutex.RLock()
	defer mutex.RUnlock()
	for uuid, conn := range conns {
		err := conn.WriteJSON(&message{Subject: subject, Payload: msg})
		if err != nil {
			debug.Debugf("failed to broadcast message to client (uuid: %s): %v", uuid, err)
		}
	}
	return nil
}
