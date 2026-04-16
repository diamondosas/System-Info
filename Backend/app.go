package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Start the HTTP API server in a separate goroutine
	go a.startHTTPServer()

	// Start network sampling
	go sampleNet(1 * time.Second)
}

// This function will be callable from JS
func (a *App) Greet(name string) string {
	return "Hello, " + name + "! 👋 from Go Backend"
}

// This function will be callable from JS to get actual specs
func (a *App) GetSpecs() (FullSpecs, error) {
	return a.collectAll()
}

// helper: convert bytes to MB (rounded)
func bytesToMB(b uint64) uint64 {
	return b / (1024 * 1024)
}

// Helper to get vendor from PNPDeviceID (e.g., VEN_10DE -> NVIDIA)
func getVendorFromPNP(pnpID string) string {
	pnpLower := strings.ToLower(pnpID)
	if strings.Contains(pnpLower, "ven_10de") {
		return "NVIDIA"
	} else if strings.Contains(pnpLower, "ven_1002") {
		return "AMD"
	} else if strings.Contains(pnpLower, "ven_8086") {
		return "Intel"
	}
	return ""
}

// Helper to infer vendor from name (fallback)
func getVendorFromName(name string) string {
	nameLower := strings.ToLower(name)
	if strings.Contains(nameLower, "nvidia") {
		return "NVIDIA"
	} else if strings.Contains(nameLower, "amd") || strings.Contains(nameLower, "radeon") {
		return "AMD"
	} else if strings.Contains(nameLower, "intel") {
		return "Intel"
	}
	return "Unknown"
}

// assemble specs
func (a *App) collectAll() (FullSpecs, error) {
	var specs FullSpecs
	// OS
	if osinfo, err := a.collectOS(); err == nil {
		specs.OS = osinfo
	}
	// CPU
	if c, err := a.collectCPU(); err == nil {
		specs.CPU = c
	}
	// Memory
	if m, err := a.collectMemory(); err == nil {
		specs.Memory = m
	}
	// GPU
	if g, err := a.collectGPU(); err == nil {
		specs.GPU = g
	}
	// Storage
	if s, err := a.collectStorage(); err == nil {
		specs.Storage = s
	}
	// Network
	if n, err := a.collectNetwork(); err == nil {
		specs.Network = n
	}
	// Battery
	if b, err := a.collectBattery(); err == nil {
		specs.Battery = b
	}
	// Sensors
	if s, err := a.collectSensors(); err == nil {
		specs.Sensors = s
	}
	// Processes
	if p, err := a.collectProcesses(); err == nil {
		specs.Processes = p
	}
	return specs, nil
}

// HTTP handlers
func (a *App) specsHandler(w http.ResponseWriter, r *http.Request) {
	specs, err := a.collectAll()
	if err != nil {
		http.Error(w, "Failed to collect specs: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(specs)
}

func (a *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

// startHTTPServer starts the HTTP API server
func (a *App) startHTTPServer() {
	// Serve the static frontend files
	fs := http.FileServer(http.Dir("../Frontend"))
	http.Handle("/", fs)

	http.HandleFunc("/api/specs", a.specsHandler)
	http.HandleFunc("/api/health", a.healthHandler)

	addr := ":9999"
	log.Printf("Starting HTTP server on %s ...", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Printf("HTTP server failed: %v", err)
	}
}
