package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/host"
)

func (a *App) collectOS() (OSInfo, error) {
	hi, err := host.Info()
	if err != nil {
		return OSInfo{}, err
	}
	u, _ := host.Uptime()

	bootTimeStr := "Unknown"
	bt, err := host.BootTime()
	if err == nil {
		bootTimeStr = time.Unix(int64(bt), 0).Format("2006-01-02 15:04:05")
	}

	return OSInfo{
		Name:         hi.OS + " " + hi.PlatformVersion,
		Version:      hi.KernelVersion,
		Architecture: hi.KernelArch,
		Uptime:       formatUptime(u),
		BootTime:     bootTimeStr,
	}, nil
}

func formatUptime(seconds uint64) string {
	d := time.Duration(seconds) * time.Second
	days := d / (24 * time.Hour)
	d -= days * 24 * time.Hour
	hours := d / time.Hour
	d -= hours * time.Hour
	mins := d / time.Minute
	return fmt.Sprintf("%d days, %d hours, %d minutes", days, hours, mins)
}
