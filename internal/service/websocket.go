package service

import (
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

var debug = debugo.New("websocket:service")

var conns = make(map[string]*websocket.Conn)

func (s *WebsocketService) AddConnection(uuid string, conn *websocket.Conn) {
	conns[uuid] = conn
}

func (s *WebsocketService) RemoveConnection(uuid string, conn *websocket.Conn) {
	delete(conns, uuid)
}

func (s *WebsocketService) Broadcast(subject Subject, msg any) error {
	for uuid, conn := range conns {
		err := conn.WriteJSON(&message{Subject: subject, Payload: msg})
		if err != nil {
			debug.Debugf("failed to broadcast message (uuid: %s): %v", uuid, err)
		}
	}
	return nil
}
