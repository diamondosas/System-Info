package main

import (
	"regexp"

	"github.com/StackExchange/wmi"
	"github.com/shirou/gopsutil/v3/disk"
)

// parse WMI association strings to get DeviceID values
func extractDeviceIDFromAssoc(s string) string {
	re := regexp.MustCompile(`DeviceID="([^"]+)"`)
	m := re.FindStringSubmatch(s)
	if len(m) >= 2 {
		return m[1]
	}
	return ""
}

func (a *App) collectStorage() ([]StorageEntry, error) {
	var wLogical []Win32_LogicalDisk
	err := wmi.Query("SELECT DeviceID, Size, FreeSpace, FileSystem, VolumeName FROM Win32_LogicalDisk WHERE DriveType=3 OR DriveType=2 OR DriveType=4", &wLogical)
	if err != nil {
		// fallback to gopsutil only
	}
	var drives []Win32_DiskDrive
	_ = wmi.Query("SELECT DeviceID, Model, Size, MediaType, Index FROM Win32_DiskDrive", &drives)
	var partitions []Win32_DiskPartition
	_ = wmi.Query("SELECT DeviceID, Name, Index, DiskIndex FROM Win32_DiskPartition", &partitions)
	var assoc []Win32_LogicalDiskToPartition
	_ = wmi.Query("SELECT Antecedent, Dependent FROM Win32_LogicalDiskToPartition", &assoc)

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
	logicalToDriveModel := map[string]Win32_DiskDrive{}

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

	entries := []StorageEntry{}
	for _, ld := range wLogical {
		capMB := uint64(0)
		freeMB := uint64(0)
		if ld.Size != nil {
			capMB = bytesToMB(*ld.Size)
		} else {
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
