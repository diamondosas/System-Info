package main

import (
	"strings"

	"github.com/StackExchange/wmi"
)

func (a *App) collectGPU() (GPUInfo, error) {
	var vcds []Win32_VideoController
	err := wmi.Query("SELECT Name, AdapterRAM, PNPDeviceID, DriverVersion FROM Win32_VideoController", &vcds)
	if err != nil {
		return GPUInfo{}, err
	}
	if len(vcds) == 0 {
		return GPUInfo{Vendor: "Unknown", Model: "Unknown", VRAMMB: 0}, nil
	}

	var realGPUs []Win32_VideoController
	for _, vc := range vcds {
		nameLower := strings.ToLower(vc.Name)
		if strings.Contains(nameLower, "basic display") ||
			strings.Contains(nameLower, "virtual display") ||
			(strings.Contains(nameLower, "microsoft") && !strings.Contains(nameLower, "nvidia") && !strings.Contains(nameLower, "amd") && !strings.Contains(nameLower, "intel")) {
			continue
		}
		realGPUs = append(realGPUs, vc)
	}

	if len(realGPUs) == 0 {
		g := vcds[0]
		vendor := getVendorFromName(g.Name)
		vram := int64(0)
		if g.AdapterRAM != nil {
			vram = int64(*g.AdapterRAM) / (1024 * 1024)
		}
		return GPUInfo{
			Vendor: vendor,
			Model:  g.Name,
			VRAMMB: vram,
		}, nil
	}

	var selected Win32_VideoController
	maxVRAM := int64(-1)
	for _, g := range realGPUs {
		vram := int64(0)
		if g.AdapterRAM != nil {
			vram = int64(*g.AdapterRAM) / (1024 * 1024)
		}
		if vram > maxVRAM {
			maxVRAM = vram
			selected = g
		}
	}

	vendor := getVendorFromPNP(selected.PNPDeviceID)
	if vendor == "" {
		vendor = getVendorFromName(selected.Name)
	}
	vram := maxVRAM

	return GPUInfo{
		Vendor: vendor,
		Model:  selected.Name,
		VRAMMB: vram,
	}, nil
}
