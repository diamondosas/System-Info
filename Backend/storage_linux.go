// +build !windows

package main

import (
	"github.com/shirou/gopsutil/v3/disk"
)

func (a *App) collectStorage() ([]StorageEntry, error) {
	entries := []StorageEntry{}
	parts, _ := disk.Partitions(true)
	for _, p := range parts {
		if usage, err := disk.Usage(p.Mountpoint); err == nil {
			entries = append(entries, StorageEntry{
				Device:   p.Mountpoint,
				Model:    "Unknown",
				Type:     p.Fstype,
				Capacity: bytesToMB(usage.Total),
				Free:     bytesToMB(usage.Free),
			})
		}
	}
	return entries, nil
}
