package discovery

import (
	"context"
	"dcp/internal/gateway"
	"log"
	"github.com/grandcat/zeroconf"
)
// internal/discovery/mdns.go

func StartDiscovery(peerID string, port int, bridge *gateway.Bridge) {
    // 1. تسجيل الخدمة
    server, err := zeroconf.Register(peerID, "_dcp-p2p._tcp", "local.", port, []string{"txtv=0"}, nil)
    if err != nil {
        log.Printf("mDNS Register Error: %v", err)
        return
    }
    defer server.Shutdown()

    // 2. إعداد المكتشف
    resolver, err := zeroconf.NewResolver(nil)
    if err != nil {
        log.Printf("mDNS Resolver Error: %v", err)
        return
    }

    // قناة استقبال النتائج
    entries := make(chan *zeroconf.ServiceEntry)

    // Goroutine لمعالجة النتائج وإرسالها للواجهة
    go func() {
        for entry := range entries {
            // تجاهل العقدة لنفسها
            if entry.Instance == peerID {
                continue
            }

            // تأكد من وجود IP متاح
            if len(entry.AddrIPv4) > 0 {
                peerData := map[string]interface{}{
                    "peer_id": entry.Instance,
                    "ip":      entry.AddrIPv4[0].String(),
                    "port":    entry.Port,
                }
                log.Printf("📡 New Peer Detected: %s", entry.Instance)
                bridge.SendToFrontend("peer_discovered", peerData)
            }
        }
    }()

    // البحث المستمر (بدون Timeout محدد حالياً للاختبار)
    ctx := context.Background()
    err = resolver.Browse(ctx, "_dcp-p2p._tcp", "local.", entries)
    if err != nil {
        log.Printf("mDNS Browse Error: %v", err)
    }

    <-ctx.Done()
}
