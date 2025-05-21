package main

import (
	// "fmt"
	"log"
	//"net"
	"net/http"
	"os"
	"system-broadcast-agent/pkg/zeroconftest"
	"time"
	// "os"
	// "os/signal"
	// "sync"
	// "syscall"
	// "system-broadcast-agent/pkg/agentserver"
	// "system-broadcast-agent/pkg/api"
	// "time"
)

func main() {
	// var agentMap sync.Map
	// hostname, _ := os.Hostname()
	// ip := getLocalIP()

	// // 1. 수신 시작
	// go agentserver.ListenAndStore(&agentMap)

	// // 2. 송신 시작
	// go func() {
	// 	for {
	// 		info := agentserver.AgentInfo{
	// 			Hostname:  hostname,
	// 			IP:        ip,
	// 			Port:      9999,
	// 			Timestamp: time.Now(),
	// 			Status:    "healthy",
	// 		}
	// 		agentserver.BroadcastMyInfo(info)
	// 		time.Sleep(5 * time.Second)
	// 	}
	// }()

	// // 2-1 : TTL 검증
	// agentserver.StartAgentTTLChecker(&agentMap, 15*time.Second)

	// // 3. API 시작
	// api.SetupAPI(&agentMap)

	// // 4. 종료 대기
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	// <-c
	// fmt.Println("Agent 종료됨")

	hostname, _ := os.Hostname()

	go func() {
		for {
			server, err := zeroconftest.RegisterZeroconfMultiIP(hostname)
			if err != nil {
				log.Println("Zeroconf 등록 실패:", err)
			} else {
				log.Println("Zeroconf 등록 완료")
			}
			time.Sleep(60 * time.Second)
			server.Shutdown()
		}
	}()

	http.HandleFunc("/discovery", zeroconftest.DiscoveryHandler)
	log.Println("API 서버 실행 중: :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

// func getLocalIP() string {
// 	interfaces, err := net.Interfaces()
// 	if err != nil {
// 		return ""
// 	}

// 	for _, i := range interfaces {
// 		addrs, err := i.Addrs()
// 		if err != nil {
// 			continue
// 		}

// 		for _, addr := range addrs {
// 			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
// 				return ipnet.IP.String()
// 			}
// 		}
// 	}
// 	return ""
// }
