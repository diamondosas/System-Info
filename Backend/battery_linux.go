// +build !windows

package main

import "fmt"

func (a *App) collectBattery() (BatteryInfo, error) {
	return BatteryInfo{
		Percentage: -1,
		Status:     fmt.Sprintf("Battery collection not supported on this platform"),
	}, nil
}
