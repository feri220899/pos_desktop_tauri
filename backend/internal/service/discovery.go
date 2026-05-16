package service

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/grandcat/zeroconf"
)

const (
	mdnsService = "_pos._tcp"
	mdnsDomain  = "local."
)

type MasterInfo struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

var (
	mu         sync.Mutex
	advertiser *zeroconf.Server
	scanCancel context.CancelFunc
)

func Advertise(name string, port int) error {
	mu.Lock()
	defer mu.Unlock()

	if advertiser != nil {
		advertiser.Shutdown()
		advertiser = nil
	}

	srv, err := zeroconf.Register(name, mdnsService, mdnsDomain, port, nil, nil)
	if err != nil {
		return err
	}
	advertiser = srv
	log.Printf("[discovery] advertising: %s on port %d", name, port)
	return nil
}

func StopAdvertise() {
	mu.Lock()
	defer mu.Unlock()
	if advertiser != nil {
		advertiser.Shutdown()
		advertiser = nil
		log.Println("[discovery] stopped advertising")
	}
}

func Scan(timeout time.Duration) ([]MasterInfo, error) {
	mu.Lock()
	if scanCancel != nil {
		scanCancel()
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	scanCancel = cancel
	mu.Unlock()
	defer func() {
		mu.Lock()
		scanCancel = nil
		mu.Unlock()
	}()

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return nil, err
	}

	entries := make(chan *zeroconf.ServiceEntry)
	var results []MasterInfo
	var resultMu sync.Mutex

	go func() {
		for entry := range entries {
			ip := pickIPv4(entry.AddrIPv4)
			if ip == "" {
				continue
			}
			resultMu.Lock()
			if !alreadyFound(results, ip, entry.Port) {
				results = append(results, MasterInfo{
					Name: entry.ServiceInstanceName(),
					IP:   ip,
					Port: entry.Port,
				})
				log.Printf("[discovery] found: %s at %s:%d", entry.ServiceInstanceName(), ip, entry.Port)
			}
			resultMu.Unlock()
		}
	}()

	if err := resolver.Browse(ctx, mdnsService, mdnsDomain, entries); err != nil {
		return nil, err
	}

	<-ctx.Done()

	resultMu.Lock()
	defer resultMu.Unlock()
	if results == nil {
		results = []MasterInfo{}
	}
	return results, nil
}

func StopScan() {
	mu.Lock()
	defer mu.Unlock()
	if scanCancel != nil {
		scanCancel()
		scanCancel = nil
	}
}

func pickIPv4(addrs []net.IP) string {
	for _, a := range addrs {
		if a.To4() != nil {
			return a.String()
		}
	}
	return ""
}

func alreadyFound(list []MasterInfo, ip string, port int) bool {
	for _, m := range list {
		if m.IP == ip && m.Port == port {
			return true
		}
	}
	return false
}
