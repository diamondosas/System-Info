package main

import (
	"strings"

	"github.com/shirou/gopsutil/v3/host"
)

func (a *App) collectSensors() (SensorsInfo, error) {
	temps, err := host.SensorsTemperatures()
	if err != nil {
		return SensorsInfo{}, err
	}
	cpuSum, gpuSum := 0.0, 0.0
	cpuCount, gpuCount := 0, 0
	for _, t := range temps {
		key := strings.ToLower(t.SensorKey)
		tempC := float64(t.Temperature) / 1000.0 // millidegrees to C
		if strings.Contains(key, "cpu") || strings.Contains(key, "thermal") {
			cpuSum += tempC
			cpuCount++
		} else if strings.Contains(key, "gpu") {
			gpuSum += tempC
			gpuCount++
		}
	}
	cpuTemp := 0.0
	if cpuCount > 0 {
		cpuTemp = cpuSum / float64(cpuCount)
	}
	gpuTemp := 0.0
	if gpuCount > 0 {
		gpuTemp = gpuSum / float64(gpuCount)
	}
	return SensorsInfo{CPUTemp: cpuTemp, GPUTemp: gpuTemp}, nil
}
