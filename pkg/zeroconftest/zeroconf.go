package zeroconftest

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/grandcat/zeroconf"
)

type AgentInfo struct {
	Hostname string   `json:"hostname"`
	IPs      []string `json:"ips"`
}

func RegisterZeroconfMultiIP(hostname string) (*zeroconf.Server, error) {
	ips := GetAllUsableIPv4()
	txt := []string{"hostname=" + hostname}
	for _, ip := range ips {
		txt = append(txt, "ip="+ip)
	}

	return zeroconf.Register(
		"agent-"+hostname,
		"_agent._tcp",
		"local.",
		9999,
		txt,
		nil,
	)
}

func DiscoveryHandler(w http.ResponseWriter, r *http.Request) {
	timeout := 3 * time.Second
	if t := r.URL.Query().Get("timeout"); t != "" {
		if parsed, err := strconv.Atoi(t); err == nil {
			timeout = time.Duration(parsed) * time.Second
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	entries := make(chan *zeroconf.ServiceEntry)
	results := []AgentInfo{}

	go func(resultsChan <-chan *zeroconf.ServiceEntry) {
		for entry := range resultsChan {
			hostname := ""
			ips := []string{}
			for _, txt := range entry.Text {
				if strings.HasPrefix(txt, "hostname=") {
					hostname = strings.TrimPrefix(txt, "hostname=")
				} else if strings.HasPrefix(txt, "ip=") {
					ips = append(ips, strings.TrimPrefix(txt, "ip="))
				}
			}
			if hostname != "" && len(ips) > 0 {
				results = append(results, AgentInfo{
					Hostname: hostname,
					IPs:      ips,
				})
			}
		}
	}(entries)

	resolver, _ := zeroconf.NewResolver(nil)
	_ = resolver.Browse(ctx, "_agent._tcp", "local.", entries)

	<-ctx.Done()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func GetAllUsableIPv4() []string {
	result := []string{}
	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		if (iface.Flags&net.FlagUp) == 0 || (iface.Flags&net.FlagLoopback) != 0 {
			continue
		}
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
				ip := ipnet.IP.String()
				if strings.HasPrefix(ip, "169.") || strings.HasPrefix(ip, "172.17.") {
					continue
				}
				result = append(result, ip)
			}
		}
	}
	return result
}
