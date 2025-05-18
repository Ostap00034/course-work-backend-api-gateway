// internal/offer/hub.go
package offer

import (
    "sync"
    "github.com/gorilla/websocket"
)

type Hub struct {
    // map[orderID]set of connections
    subs   map[string]map[*websocket.Conn]bool
    mu     sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        subs: make(map[string]map[*websocket.Conn]bool),
    }
}

func (h *Hub) Subscribe(orderID string, conn *websocket.Conn) {
    h.mu.Lock()
    defer h.mu.Unlock()
    if h.subs[orderID] == nil {
        h.subs[orderID] = make(map[*websocket.Conn]bool)
    }
    h.subs[orderID][conn] = true
}

func (h *Hub) Unsubscribe(orderID string, conn *websocket.Conn) {
    h.mu.Lock()
    defer h.mu.Unlock()
    if conns := h.subs[orderID]; conns != nil {
        delete(conns, conn)
        if len(conns) == 0 {
            delete(h.subs, orderID)
        }
    }
}

func (h *Hub) Broadcast(orderID string, message interface{}) {
    h.mu.RLock()
    defer h.mu.RUnlock()
    for conn := range h.subs[orderID] {
        conn.WriteJSON(message)
    }
}