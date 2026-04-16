// +build !windows

package main

func (a *App) collectGPU() (GPUInfo, error) {
	return GPUInfo{
		Vendor: "Unknown",
		Model:  "GPU collection not supported on this platform",
		VRAMMB: 0,
	}, nil
}
