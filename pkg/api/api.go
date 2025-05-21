package api

import (
	"encoding/json"
	"net/http"
	"sync"
	"system-broadcast-agent/pkg/agentserver"
)

func SetupAPI(agentMap *sync.Map) {
	http.HandleFunc("/agents", func(w http.ResponseWriter, r *http.Request) {
		agents := []agentserver.AgentInfo{}
		agentMap.Range(func(_, v interface{}) bool {
			agents = append(agents, v.(agentserver.AgentInfo))
			return true
		})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(agents)
	})
	go http.ListenAndServe(":8080", nil)
}
