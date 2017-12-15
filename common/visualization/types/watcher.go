package types

import (
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type Watcher struct {
	id   string
	conn *websocket.Conn
}

func NewWatcher(conn *websocket.Conn) *Watcher {
	return &Watcher{
		id:   uuid.NewV4().String(),
		conn: conn,
	}
}

func (watcher *Watcher) GetId() string {
	return watcher.id
}

func (watcher *Watcher) GetConn() *websocket.Conn {
	return watcher.conn
}
