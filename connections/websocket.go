package connections

import (
    "sync"
    "github.com/gorilla/websocket"
)

var Clients sync.Map

func GetClients() *sync.Map {
    return &Clients
}

func StoreClient(userID string, conn *websocket.Conn) {
    Clients.Store(userID, conn)
}

func DeleteClient(userID string) {
    Clients.Delete(userID)
}

func GetClient(userID string) (*websocket.Conn, bool) {
    if conn, ok := Clients.Load(userID); ok {
        return conn.(*websocket.Conn), true
    }
    return nil, false
}