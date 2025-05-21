package agentserver

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/grandcat/zeroconf"
)

type AgentInfo struct {
	Hostname  string    `json:"hostname"`
	IP        string    `json:"ip"`
	Port      int       `json:"port"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

func BroadcastMyInfo(info AgentInfo) {
	data, _ := json.Marshal(info)
	encoded := base64.StdEncoding.EncodeToString(data)
	txt := []string{"payload=" + encoded}

	server, err := zeroconf.Register(
		info.Hostname, "_agent._udp", "local.", info.Port, txt, nil,
	)
	if err == nil {
		time.Sleep(2 * time.Second)
		server.Shutdown()
	}
}

func ListenAndStore(agentMap *sync.Map) {
	resolver, _ := zeroconf.NewResolver(nil)
	entries := make(chan *zeroconf.ServiceEntry)

	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			for _, txt := range entry.Text {
				if strings.HasPrefix(txt, "payload=") {
					decoded, _ := base64.StdEncoding.DecodeString(strings.TrimPrefix(txt, "payload="))
					var info AgentInfo
					if json.Unmarshal(decoded, &info) == nil {
						agentMap.Store(info.Hostname, info)
					}
				}
			}
		}
	}(entries)

	for {
		ctx := context.Background()
		resolver.Browse(ctx, "_agent._udp", "local.", entries)
		time.Sleep(10 * time.Second)
	}
}

func StartAgentTTLChecker(agentMap *sync.Map, timeout time.Duration) {
	go func() {
		for {
			now := time.Now()
			agentMap.Range(func(key, value interface{}) bool {
				info := value.(AgentInfo)
				if now.Sub(info.Timestamp) > timeout {
					log.Printf("TTL 만료: %s 삭제됨", info.Hostname)
					agentMap.Delete(key)
				}
				return true
			})
			time.Sleep(5 * time.Second) // 주기적 검사
		}
	}()
}
