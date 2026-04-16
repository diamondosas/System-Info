package main

type Win32_OperatingSystem struct {
	Name           string
	Version        string
	BuildNumber    string
	OSArchitecture string
	LastBootUpTime string
}

type Win32_VideoController struct {
	Name          string
	AdapterRAM    *uint64
	PNPDeviceID   string
	DriverVersion string
}

type Win32_DiskDrive struct {
	DeviceID  string
	Model     string
	Size      *uint64
	MediaType *string
	Index     *uint32
}

type Win32_LogicalDisk struct {
	DeviceID   string
	Size       *uint64
	FreeSpace  *uint64
	FileSystem *string
	VolumeName *string
}

type Win32_DiskPartition struct {
	DeviceID  string
	Name      string
	Index     *uint32
	DiskIndex *uint32
}

type Win32_LogicalDiskToPartition struct {
	Antecedent string
	Dependent  string
}

type Win32_Battery struct {
	EstimatedChargeRemaining *uint16
	Status                   *uint16
}

type OSInfo struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Architecture string `json:"architecture"`
	Uptime       string `json:"uptime"`
	BootTime     string `json:"boot_time"`
}

type CPUInfo struct {
	Model        string  `json:"model"`
	Cores        int     `json:"cores"`
	Threads      int     `json:"threads"`
	FrequencyMHz int64   `json:"frequency"`
	UsagePercent float64 `json:"usage_percent"`
}

type MemoryInfo struct {
	TotalMB        uint64 `json:"total"`
	AvailableMB    uint64 `json:"available"`
	UsedMB         uint64 `json:"used"`
	SwapTotalMB    uint64 `json:"swap_total"`
	SwapFreeMB     uint64 `json:"swap_free"`
	SwapUsedMB     uint64 `json:"swap_used"`
	SwapUsagePercent float64 `json:"swap_usage_percent"`
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
	Capacity uint64 `json:"capacity"`
	Free     uint64 `json:"free"`
}

type NetworkEntry struct {
	Interface     string `json:"interface"`
	IPAddress     string `json:"ip_address"`
	MACAddress    string `json:"mac_address"`
	BandwidthDown uint64 `json:"bandwidth_down"`
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
