package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/StackExchange/wmi"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// This function will be callable from JS
func (a *App) Greet(name string) string {
    return "Hello, " + name + "! 👋 from Go Backend"
}
// --- WMI structs (fields we care about) ---
type Win32_OperatingSystem struct {
	Name           string
	Version        string
	BuildNumber    string
	OSArchitecture string
	LastBootUpTime string
}

type Win32_VideoController struct {
	Name        string
	AdapterRAM  *uint64 // bytes; pointer because sometimes nil
	PNPDeviceID string
	DriverVersion string
}

type Win32_DiskDrive struct {
	DeviceID   string // e.g. \\.\PHYSICALDRIVE0
	Model      string
	Size       *uint64 // bytes
	MediaType  *string
	Index      *uint32
}

type Win32_LogicalDisk struct {
	DeviceID    string // "C:"
	Size        *uint64
	FreeSpace   *uint64
	FileSystem  *string
	VolumeName  *string
}

type Win32_DiskPartition struct {
	DeviceID  string // "Disk #0, Partition #0"
	Name      string
	Index     *uint32
	DiskIndex *uint32
}

type Win32_LogicalDiskToPartition struct {
	Antecedent string // gives partition reference
	Dependent  string // gives logical disk reference
}

type Win32_Battery struct {
	EstimatedChargeRemaining *uint16
	Status                   *uint16
}

// --- JSON response types (match your example) ---
type OSInfo struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Architecture string `json:"architecture"`
	Uptime       string `json:"uptime"`
}

type CPUInfo struct {
	Model        string  `json:"model"`
	Cores        int     `json:"cores"`
	Threads      int     `json:"threads"`
	FrequencyMHz int64   `json:"frequency"`
	UsagePercent float64 `json:"usage_percent"`
}

type MemoryInfo struct {
	TotalMB     uint64 `json:"total"`
	AvailableMB uint64 `json:"available"`
	UsedMB      uint64 `json:"used"`
}

type GPUInfo struct {
	Vendor string `json:"vendor"`
	Model  string `json:"model"`
	VRAMMB int64  `json:"vram"`
}

type StorageEntry struct {
	Device   string `json:"device"`
	Model    string `json:"model"`
	Type     string `json:"type"`
	Capacity uint64 `json:"capacity"` // MB
	Free     uint64 `json:"free"`     // MB
}

type NetworkEntry struct {
	Interface     string `json:"interface"`
	IPAddress     string `json:"ip_address"`
	MACAddress    string `json:"mac_address"`
	BandwidthDown uint64 `json:"bandwidth_down"` // bps
}

type BatteryInfo struct {
	Percentage int    `json:"percentage"`
	Status     string `json:"status"`
}

type SensorsInfo struct {
	CPUTemp float64 `json:"cpu_temp"`
	GPUTemp float64 `json:"gpu_temp"`
}

type ProcessEntry struct {
	PID      int32   `json:"pid"`
	Name     string  `json:"name"`
	MemoryMB float32 `json:"memory_mb"`
}

type FullSpecs struct {
	OS        OSInfo         `json:"os"`
	CPU       CPUInfo        `json:"cpu"`
	Memory    MemoryInfo     `json:"memory"`
	GPU       GPUInfo        `json:"gpu"`
	Storage   []StorageEntry `json:"storage"`
	Network   []NetworkEntry `json:"network"`
	Battery   BatteryInfo    `json:"battery"`
	Sensors   SensorsInfo    `json:"sensors"`
	Processes []ProcessEntry `json:"processes"`
}

// --- Globals for network sampling ---
var (
	netLock sync.RWMutex
	// store bits per second per interface name
	netBPS = map[string]uint64{}
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	
	// Start the HTTP API server in a separate goroutine
	go a.startHTTPServer()
	
	// Start network sampling
	go sampleNet(1 * time.Second)
	
	// Show message to user
}

// sampleNet periodically samples net.IOCounters to compute bps per interface.
// runs as a background goroutine when server starts.
func sampleNet(interval time.Duration) {
	prev, _ := net.IOCounters(true)
	prevMap := make(map[string]net.IOCountersStat)
	for _, s := range prev {
		prevMap[s.Name] = s
	}
	ticker := time.NewTicker(interval)
	for range ticker.C {
		cur, err := net.IOCounters(true)
		if err != nil {
			continue
		}
		nowMap := make(map[string]net.IOCountersStat)
		for _, s := range cur {
			nowMap[s.Name] = s
		}
		netLock.Lock()
		for name, curStat := range nowMap {
			if p, ok := prevMap[name]; ok {
				// bytes received delta / seconds
				deltaBytes := curStat.BytesRecv - p.BytesRecv
				bps := (deltaBytes * 8) / uint64(interval.Seconds())
				netBPS[name] = bps
			} else {
				netBPS[name] = 0
			}
		}
		netLock.Unlock()
		prevMap = nowMap
	}
}

// helper: convert bytes to MB (rounded)
func bytesToMB(b uint64) uint64 {
	return b / (1024 * 1024)
}

// parse WMI association strings to get DeviceID values
// e.g. Antecedent: Win32_DiskPartition.DeviceID="Disk #0, Partition #0"
func extractDeviceIDFromAssoc(s string) string {
	// find first occurrence of DeviceID="..." and return content
	re := regexp.MustCompile(`DeviceID="([^"]+)"`)
	m := re.FindStringSubmatch(s)
	if len(m) >= 2 {
		return m[1]
	}
	return ""
}

// --- Collectors ---
func (a *App) collectOS() (OSInfo, error) {
	var osInfos []Win32_OperatingSystem
	err := wmi.Query("SELECT Name, Version, OSArchitecture, LastBootUpTime FROM Win32_OperatingSystem", &osInfos)
	if err != nil {
		return OSInfo{}, err
	}
	var out OSInfo
	if len(osInfos) > 0 {
		wo := osInfos[0]
		// Name often has additional text; example: "Microsoft Windows 10 Pro|C:\WINDOWS|\Device\Harddisk0\Partition4"
		// We'll strip down to first part up to the '|' or full string.
		name := wo.Name
		if idx := strings.Index(name, "|"); idx != -1 {
			name = name[:idx]
		}
		out.Name = strings.TrimSpace(name)
		out.Version = wo.Version
		out.Architecture = wo.OSArchitecture
		// compute uptime using host package as well (easier)
		u, _ := host.Uptime()
		out.Uptime = formatUptime(u)
	}
	return out, nil
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
		// cores: physical + threads
		// gopsutil's cpu.Info returns logical CPUs; get physical cores via cpu.Counts?
		phys, _ := cpu.Counts(false)
		logical, _ := cpu.Counts(true)
		info.Cores = phys
		info.Threads = logical
	}
	// usage percent (instant, interval sample 500ms)
	percent, err := cpu.Percent(500*time.Millisecond, false)
	if err == nil && len(percent) > 0 {
		info.UsagePercent = percent[0]
	}
	// if frequency absent, try to get from info entries
	if info.FrequencyMHz == 0 && len(ci) > 0 {
		info.FrequencyMHz = int64(ci[0].Mhz)
	}
	return info, nil
}

func (a *App) collectMemory() (MemoryInfo, error) {
	m, err := mem.VirtualMemory()
	if err != nil {
		return MemoryInfo{}, err
	}
	return MemoryInfo{
		TotalMB:     bytesToMB(m.Total),
		AvailableMB: bytesToMB(m.Available),
		UsedMB:      bytesToMB(m.Total - m.Available),
	}, nil
}

func (a *App) collectGPU() (GPUInfo, error) {
	var vcds []Win32_VideoController
	err := wmi.Query("SELECT Name, AdapterRAM, PNPDeviceID, DriverVersion FROM Win32_VideoController", &vcds)
	if err != nil {
		return GPUInfo{}, err
	}
	if len(vcds) == 0 {
		return GPUInfo{Vendor: "Unknown", Model: "Unknown", VRAMMB: 0}, nil
	}
	// pick first GPU (primary)
	g := vcds[0]
	vendor := "Unknown"
	// try infer vendor from name
	if strings.Contains(strings.ToLower(g.Name), "nvidia") {
		vendor = "NVIDIA"
	} else if strings.Contains(strings.ToLower(g.Name), "amd") || strings.Contains(strings.ToLower(g.Name), "radeon") {
		vendor = "AMD"
	} else if strings.Contains(strings.ToLower(g.Name), "intel") {
		vendor = "Intel"
	}
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

func (a *App) collectStorage() ([]StorageEntry, error) {
	// Get logical disks via WMI & gopsutil
	var wLogical []Win32_LogicalDisk
	err := wmi.Query("SELECT DeviceID, Size, FreeSpace, FileSystem, VolumeName FROM Win32_LogicalDisk WHERE DriveType=3 OR DriveType=2 OR DriveType=4", &wLogical)
	if err != nil {
		// fallback to gopsutil only
	}
	// Get physical drives
	var drives []Win32_DiskDrive
	_ = wmi.Query("SELECT DeviceID, Model, Size, MediaType, Index FROM Win32_DiskDrive", &drives)
	// Get partitions + associations to map to logical disks
	var partitions []Win32_DiskPartition
	_ = wmi.Query("SELECT DeviceID, Name, Index, DiskIndex FROM Win32_DiskPartition", &partitions)
	var assoc []Win32_LogicalDiskToPartition
	_ = wmi.Query("SELECT Antecedent, Dependent FROM Win32_LogicalDiskToPartition", &assoc)

	// maps to help matching
	diskIndexToDrive := map[uint32]Win32_DiskDrive{}
	for _, d := range drives {
		if d.Index != nil {
			diskIndexToDrive[*d.Index] = d
		}
	}
	partitionNameToPartition := map[string]Win32_DiskPartition{}
	for _, p := range partitions {
		partitionNameToPartition[p.DeviceID] = p
	}
	logicalToDriveModel := map[string]Win32_DiskDrive{} // logicalDeviceID -> disk drive (best-effort)

	for _, a := range assoc {
		partDev := extractDeviceIDFromAssoc(a.Antecedent)
		logDev := extractDeviceIDFromAssoc(a.Dependent)
		if partDev == "" || logDev == "" {
			continue
		}
		if p, ok := partitionNameToPartition[partDev]; ok && p.DiskIndex != nil {
			if d, ok2 := diskIndexToDrive[*p.DiskIndex]; ok2 {
				logicalToDriveModel[logDev] = d
			}
		}
	}

	// Build storage entries - prefer WMI logical disk sizes if available; otherwise use gopsutil
	entries := []StorageEntry{}
	for _, ld := range wLogical {
		capMB := uint64(0)
		freeMB := uint64(0)
		if ld.Size != nil {
			capMB = bytesToMB(*ld.Size)
		} else {
			// fallback: use gopsutil disk.Usage
			if usage, err := disk.Usage(ld.DeviceID + "\\"); err == nil {
				capMB = bytesToMB(usage.Total)
				freeMB = bytesToMB(usage.Free)
			}
		}
		if ld.FreeSpace != nil {
			freeMB = bytesToMB(*ld.FreeSpace)
		}
		model := "Unknown"
		mediaType := "Unknown"
		if drive, ok := logicalToDriveModel[ld.DeviceID]; ok {
			if drive.Model != "" {
				model = drive.Model
			}
			if drive.MediaType != nil {
				mediaType = *drive.MediaType
			}
		}
		entries = append(entries, StorageEntry{
			Device:   ld.DeviceID,
			Model:    model,
			Type:     mediaType,
			Capacity: capMB,
			Free:     freeMB,
		})
	}

	// If we didn't find anything via WMI logical, fallback to partitions from gopsutil
	if len(entries) == 0 {
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
	}
	return entries, nil
}

func (a *App) collectNetwork() ([]NetworkEntry, error) {
	// Get addrs and net interfaces
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	out := []NetworkEntry{}
	netLock.RLock()
	defer netLock.RUnlock()
	for _, inf := range ifaces {
		ip := ""
		if len(inf.Addrs) > 0 {
			ip = inf.Addrs[0].Addr
			// strip mask if present
			if idx := strings.Index(ip, "/"); idx != -1 {
				ip = ip[:idx]
			}
		}
		mac := inf.HardwareAddr
		bps := uint64(0)
		if val, ok := netBPS[inf.Name]; ok {
			bps = val
		}
		out = append(out, NetworkEntry{
			Interface:     inf.Name,
			IPAddress:     ip,
			MACAddress:    mac,
			BandwidthDown: bps,
		})
	}
	return out, nil
}

func (a *App) collectBattery() (BatteryInfo, error) {
	var bats []Win32_Battery
	_ = wmi.Query("SELECT EstimatedChargeRemaining, Status FROM Win32_Battery", &bats)
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
			status = "The device is not supported"
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
		default:
			status = fmt.Sprintf("StatusCode:%d", *b.Status)
		}
	}
	return BatteryInfo{
		Percentage: perc,
		Status:     status,
	}, nil
}

func (a *App) collectSensors() (SensorsInfo, error) {
	// Windows often doesn't expose sensors via standard gopsutil on all machines.
	// For sample code we'll return zeros or try MSAcpi_ThermalZoneTemperature WMI class.
	type MSAcpi_Thermal struct {
		CurrentTemperature *int32 // in tenths of Kelvin
		InstanceName       *string
	}
	var temps []MSAcpi_Thermal
	_ = wmi.Query("SELECT CurrentTemperature, InstanceName FROM MSAcpi_ThermalZoneTemperature", &temps)
	cpuTemp := 0.0
	if len(temps) > 0 && temps[0].CurrentTemperature != nil {
		// CurrentTemperature is in tenths of Kelvin according to some docs: (K * 10)
		k := float64(*temps[0].CurrentTemperature) / 10.0
		c := k - 273.15
		cpuTemp = c
	}
	// GPU temp via WMI is not reliable; set to 0 if not found
	return SensorsInfo{CPUTemp: cpuTemp, GPUTemp: 0.0}, nil
}

func (a *App) collectProcesses() ([]ProcessEntry, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}
	out := []ProcessEntry{}
	limit := 200 // limit how many processes to include to avoid heavy responses; adjust or remove as needed
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
	http.HandleFunc("/api/specs", a.specsHandler)
	http.HandleFunc("/api/health", a.healthHandler)

	addr := ":9999"
	log.Printf("Starting HTTP server on %s ...", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Printf("HTTP server failed: %v", err)
	}
}