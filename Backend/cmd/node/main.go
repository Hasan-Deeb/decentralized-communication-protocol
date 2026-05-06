package main

import (
	"dcp/internal/discovery"
	"dcp/internal/gateway"
	"dcp/internal/peer"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	// تحديد المنفذ من سطر الأوامر
	port := flag.Int("port", 8080, "Main port for the node")
	flag.Parse()

	id, _ := peer.GenerateIdentity()
	myShortID := id.PeerID[:12]
	bridge := gateway.NewBridge()

	// 1. تنظيف العقد المنقطعة كل 10 ثوانٍ
	bridge.StartJanitor(10 * time.Second)

	// 2. تشغيل خادم P2P لاستقبال الـ Heartbeats والـ Pings
	p2pPort := *port + 1000
	go func() {
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", p2pPort))
		if err != nil {
			log.Fatalf("P2P Listener Error: %v", err)
		}
		for {
			conn, err := l.Accept()
			if err != nil {
				continue
			}
			go handleP2PTraffic(conn, bridge)
		}
	}()

	// 3. حلقة النبض (Heartbeat Loop) لإعلام الأقران أننا "أحياء"
	go func() {
		for {
			time.Sleep(5 * time.Second)
			// نستخدم الدالة الآمنة بدلاً من الوصول لـ mu مباشرة
			peers := bridge.GetAllPeers()
			for _, p := range peers {
				go func(targetIP string, targetPort int) {
					conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", targetIP, targetPort+1000), 1*time.Second)
					if err != nil {
						return
					}
					defer conn.Close()
					json.NewEncoder(conn).Encode(map[string]interface{}{
						"type": "heartbeat",
						"id":   myShortID,
						"port": *port,
						"ip":   "127.0.0.1",
					})
				}(p.IP, p.Port)
			}
		}
	}()

	// 4. تفعيل اكتشاف mDNS
	go discovery.StartDiscovery(myShortID, *port, bridge)

	// 5. مسارات الويب للواجهة والـ WebSocket
	http.HandleFunc("/ws", bridge.HandleWS)
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, `{"peer_id": "%s", "port": %d}`, id.PeerID, *port)
	})

	log.Printf("🚀 Node [%s] started on port %d (P2P: %d)", myShortID, *port, p2pPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func handleP2PTraffic(conn net.Conn, b *gateway.Bridge) {
	defer conn.Close()
	var msg map[string]interface{}
	if err := json.NewDecoder(conn).Decode(&msg); err != nil {
		return
	}

	msgType, ok := msg["type"].(string)
	if !ok {
		return
	}

	switch msgType {
	case "heartbeat":
		senderID := msg["id"].(string)
		senderPort := int(msg["port"].(float64))
		senderIP := msg["ip"].(string)
		b.UpdatePeer(senderID, senderIP, senderPort)

	case "ping":
		b.SendToFrontend("ping_received", map[string]string{"status": "received"})
	}
}
