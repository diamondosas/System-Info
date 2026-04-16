package main

import (
	"fmt"

	"github.com/StackExchange/wmi"
)

func (a *App) collectBattery() (BatteryInfo, error) {
	var bats []Win32_Battery
	err := wmi.Query("SELECT EstimatedChargeRemaining, Status FROM Win32_Battery", &bats)
	if err != nil {
		return BatteryInfo{Percentage: -1, Status: "Error"}, err
	}
	if len(bats) == 0 {
		return BatteryInfo{Percentage: -1, Status: "NoBattery"}, nil
	}
	b := bats[0]
	perc := -1
	if b.EstimatedChargeRemaining != nil {
		perc = int(*b.EstimatedChargeRemaining)
	}
	status := "Unknown"
	if b.Status != nil {
		switch *b.Status {
		case 1:
			status = "Not supported"
		case 2:
			status = "Discharging"
		case 3:
			status = "Fully Charged"
		case 4:
			status = "Low"
		case 5:
			status = "Critical"
		case 6:
			status = "Charging"
		case 7:
			status = "Charging, high"
		case 8:
			status = "Charging, low"
		case 9:
			status = "Charging, critical"
		case 10:
			status = "Undefined"
		case 11:
			status = "Partially Charged"
		default:
			status = fmt.Sprintf("StatusCode:%d", *b.Status)
		}
	}
	return BatteryInfo{
		Percentage: perc,
		Status:     status,
	}, nil
}
