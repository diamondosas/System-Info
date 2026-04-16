package main

import (
	"github.com/shirou/gopsutil/v3/mem"
)

func (a *App) collectMemory() (MemoryInfo, error) {
	m, err := mem.VirtualMemory()
	if err != nil {
		return MemoryInfo{}, err
	}

	s, err := mem.SwapMemory()
	var swapTotal, swapFree, swapUsed uint64
	var swapUsagePercent float64
	if err == nil {
		swapTotal = bytesToMB(s.Total)
		swapFree = bytesToMB(s.Free)
		swapUsed = bytesToMB(s.Used)
		swapUsagePercent = s.UsedPercent
	}

	return MemoryInfo{
		TotalMB:          bytesToMB(m.Total),
		AvailableMB:      bytesToMB(m.Available),
		UsedMB:           bytesToMB(m.Total - m.Available),
		SwapTotalMB:      swapTotal,
		SwapFreeMB:       swapFree,
		SwapUsedMB:       swapUsed,
		SwapUsagePercent: swapUsagePercent,
	}, nil
}
