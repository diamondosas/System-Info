package main

import (
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/net"
)

var (
	netLock sync.RWMutex
	netBPS  = map[string]uint64{}
)

func sampleNet(interval time.Duration) {
	prev, _ := net.IOCounters(true)
	prevMap := make(map[string]net.IOCountersStat)
	for _, s := range prev {
		prevMap[s.Name] = s
	}
	ticker := time.NewTicker(interval)
	for range ticker.C {
		cur, err := net.IOCounters(true)
		if err != nil {
			continue
		}
		nowMap := make(map[string]net.IOCountersStat)
		for _, s := range cur {
			nowMap[s.Name] = s
		}
		netLock.Lock()
		for name, curStat := range nowMap {
			if p, ok := prevMap[name]; ok {
				deltaBytes := curStat.BytesRecv - p.BytesRecv
				bps := (deltaBytes * 8) / uint64(interval.Seconds())
				netBPS[name] = bps
			} else {
				netBPS[name] = 0
			}
		}
		netLock.Unlock()
		prevMap = nowMap
	}
}

func (a *App) collectNetwork() ([]NetworkEntry, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	out := []NetworkEntry{}
	netLock.RLock()
	defer netLock.RUnlock()
	for _, inf := range ifaces {
		ip := ""
		if len(inf.Addrs) > 0 {
			ip = inf.Addrs[0].Addr
			if idx := strings.Index(ip, "/"); idx != -1 {
				ip = ip[:idx]
			}
		}
		mac := inf.HardwareAddr
		bps := uint64(0)
		if val, ok := netBPS[inf.Name]; ok {
			bps = val
		}
		out = append(out, NetworkEntry{
			Interface:     inf.Name,
			IPAddress:     ip,
			MACAddress:    mac,
			BandwidthDown: bps,
		})
	}
	return out, nil
}
