package main

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

func (a *App) collectCPU() (CPUInfo, error) {
	info := CPUInfo{}
	ci, err := cpu.Info()
	if err != nil {
		return info, err
	}
	if len(ci) > 0 {
		info.Model = ci[0].ModelName
		if ci[0].Mhz > 0 {
			info.FrequencyMHz = int64(ci[0].Mhz)
		}
		phys, _ := cpu.Counts(false)
		logical, _ := cpu.Counts(true)
		info.Cores = phys
		info.Threads = logical
	}
	percent, err := cpu.Percent(500*time.Millisecond, false)
	if err == nil && len(percent) > 0 {
		info.UsagePercent = percent[0]
	}
	if info.FrequencyMHz == 0 && len(ci) > 0 {
		info.FrequencyMHz = int64(ci[0].Mhz)
	}
	return info, nil
}
