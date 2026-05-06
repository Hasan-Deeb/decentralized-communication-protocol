package gateway

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type PeerInfo struct {
	PeerID   string    `json:"peer_id"`
	IP       string    `json:"ip"`
	Port     int       `json:"port"`
	LastSeen time.Time `json:"-"`
}

type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Bridge struct {
	Clients    map[*websocket.Conn]bool
	PeersStore map[string]PeerInfo
	mu         sync.RWMutex // حرف صغير تعني أنها Private
}

func NewBridge() *Bridge {
	return &Bridge{
		Clients:    make(map[*websocket.Conn]bool),
		PeersStore: make(map[string]PeerInfo),
	}
}

// دالة حماية لتحويل الأنواع القادمة من JSON أو mDNS
func interfaceToInt(val interface{}) int {
	switch v := val.(type) {
	case int:
		return v
	case float64:
		return int(v)
	default:
		return 0
	}
}

// GetAllPeers تعيد نسخة من الأقران لاستخدامها في main.go بأمان
func (b *Bridge) GetAllPeers() []PeerInfo {
	b.mu.RLock()
	defer b.mu.RUnlock()
	peers := make([]PeerInfo, 0, len(b.PeersStore))
	for _, p := range b.PeersStore {
		peers = append(peers, p)
	}
	return peers
}

func (b *Bridge) StartJanitor(timeout time.Duration) {
	go func() {
		for {
			time.Sleep(2 * time.Second)
			b.mu.Lock()
			for id, peer := range b.PeersStore {
				if time.Since(peer.LastSeen) > timeout {
					log.Printf("🗑️ Node %s timed out", id)
					delete(b.PeersStore, id)
					b.broadcastToFrontend("peer_removed", map[string]string{"peer_id": id})
				}
			}
			b.mu.Unlock()
		}
	}()
}

func (b *Bridge) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	b.mu.Lock()
	b.Clients[conn] = true
	for _, peer := range b.PeersStore {
		conn.WriteJSON(WSMessage{Type: "peer_discovered", Data: peer})
	}
	b.mu.Unlock()

	go func() {
		defer func() {
			conn.Close()
			b.mu.Lock()
			delete(b.Clients, conn)
			b.mu.Unlock()
		}()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}
			var msg WSMessage
			if err := json.Unmarshal(message, &msg); err == nil && msg.Type == "send_ping" {
				b.handlePingRequest(msg.Data)
			}
		}
	}()
}

func (b *Bridge) handlePingRequest(data interface{}) {
	d := data.(map[string]interface{})
	peerIP := d["ip"].(string)
	peerPort := interfaceToInt(d["port"]) + 1000
	targetID := d["peer_id"].(string)

	go func() {
		p2pConn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", peerIP, peerPort), 2*time.Second)
		if err != nil {
			b.broadcastToFrontend("ping_sent_status", map[string]string{"target": targetID, "status": "Failed"})
			return
		}
		defer p2pConn.Close()
		json.NewEncoder(p2pConn).Encode(map[string]interface{}{
			"type": "ping", 
			"id": "External-Node",
		})
		b.broadcastToFrontend("ping_sent_status", map[string]string{"target": targetID, "status": "Delivered"})
	}()
}

// دالة التوافق مع ملف mdns.go
func (b *Bridge) SendToFrontend(msgType string, data interface{}) {
	if msgType == "peer_discovered" {
		p := data.(map[string]interface{})
		b.UpdatePeer(p["peer_id"].(string), p["ip"].(string), interfaceToInt(p["port"]))
	} else {
		b.broadcastToFrontend(msgType, data)
	}
}

func (b *Bridge) UpdatePeer(id string, ip string, port int) {
	b.mu.Lock()
	_, exists := b.PeersStore[id]
	b.PeersStore[id] = PeerInfo{
		PeerID:   id,
		IP:       ip,
		Port:     port,
		LastSeen: time.Now(),
	}
	b.mu.Unlock()

	if !exists {
		b.broadcastToFrontend("peer_discovered", b.PeersStore[id])
	}
}

func (b *Bridge) broadcastToFrontend(msgType string, data interface{}) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	msg := WSMessage{Type: msgType, Data: data}
	for client := range b.Clients {
		client.WriteJSON(msg)
	}
}
