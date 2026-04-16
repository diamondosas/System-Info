package main

import (
	"github.com/shirou/gopsutil/v3/process"
)

func (a *App) collectProcesses() ([]ProcessEntry, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}
	out := []ProcessEntry{}
	limit := 200 // limit how many processes to include to avoid heavy responses
	i := 0
	for _, p := range procs {
		if i >= limit {
			break
		}
		name, _ := p.Name()
		memInfo, _ := p.MemoryInfo()
		if memInfo == nil {
			continue
		}
		out = append(out, ProcessEntry{
			PID:      p.Pid,
			Name:     name,
			MemoryMB: float32(bytesToMB(memInfo.RSS)),
		})
		i++
	}
	return out, nil
}
