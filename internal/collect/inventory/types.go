package inventory

import (
	"fmt"
	"strconv"
	"time"

	"github.com/vignemail1/memscope/internal/model"
)

// SystemInventory represents complete system hardware inventory
type SystemInventory struct {
	Timestamp time.Time
	System    SystemInfo
	CPU       CPUInfo
	BIOS      BIOSInfo
	Memory    []MemoryModuleInfo
}

// SystemInfo represents basic system information
type SystemInfo struct {
	Manufacturer string
	Model        string
	SerialNumber string
}

// CPUInfo represents processor information
type CPUInfo struct {
	Name         string
	Manufacturer string
	Architecture string
	Cores        uint32
	Threads      uint32
	MaxClockMHz  uint32
}

// BIOSInfo represents BIOS/UEFI information
type BIOSInfo struct {
	Vendor  string
	Version string
	Date    string
}

// MemoryModuleInfo represents a physical memory module
type MemoryModuleInfo struct {
	DeviceLocator  string
	BankLabel      string
	Capacity       uint64
	Speed          uint32
	Manufacturer   string
	PartNumber     string
	SerialNumber   string
	DataWidth      uint16
	TotalWidth     uint16
	FormFactor     uint16
	MemoryType     uint16
	SPD            *SPDInfo
}

// SPDInfo contains Serial Presence Detect data
type SPDInfo struct {
	Raw     []byte
	Profile *MemoryProfile
}

// MemoryProfile represents parsed memory profile data
type MemoryProfile struct {
	Type                string
	Manufacturer        string
	PartNumber          string
	ModuleSize          uint64
	OrganizationWidth   uint16
	BankCount           uint16
	DensityPerChip      uint64
	RanksPerModule      uint16
	PrimaryBusWidth     uint16
	ExtensionBusWidth   uint16
	VoltageLevel        float64
	CycleTime           uint32
	AccessTime          uint32
	TimingParameters    TimingParameters
	SerialNumber        string
	WeekCode            uint16
	YearCode            uint16
}

// TimingParameters contains memory timing information
type TimingParameters struct {
	CL   uint16 // CAS Latency
	RCD  uint16 // RAS to CAS delay
	RP   uint16 // Row precharge time
	RAS  uint16 // Active to precharge delay
	RC   uint16 // Active to active/refresh delay
	RFC  uint16 // Refresh cycle time
	WR   uint16 // Write recovery time
	WTR  uint16 // Write to read delay
	RTP  uint16 // Read to precharge delay
	FAW  uint16 // Four activate window
}

// ToSnapshot converts SystemInventory to model.Snapshot
func (inv *SystemInventory) ToSnapshot() (*model.Snapshot, error) {
	snapshot := &model.Snapshot{
		SchemaVersion:        "1.0.0",
		ToolVersion:          "memscope-dev",
		ID:                   fmt.Sprintf("inventory-%d", inv.Timestamp.Unix()),
		CollectionStartedAt:  inv.Timestamp,
		CollectionFinishedAt: inv.Timestamp,
		Platform: model.Platform{
			OS:   "windows",
			Arch: "amd64",
		},
		Devices:      []model.Device{},
		Observations: []model.Observation{},
		Profiles:     []model.Profile{},
		Diagnostics:  []model.Diagnostic{},
		Capabilities: []model.Capability{},
	}

	// Add system device
	systemDevice := model.Device{
		ID:   "system-0",
		Kind: "system",
	}
	snapshot.Devices = append(snapshot.Devices, systemDevice)

	// Add system observations
	if inv.System.Manufacturer != "" {
		snapshot.Observations = append(snapshot.Observations, model.Observation{
			DeviceID:   "system-0",
			Scope:      "hardware",
			Parameter:  "manufacturer",
			Value:      &model.Value{Text: &inv.System.Manufacturer},
			Source:     model.WMI,
			Status:     model.Observed,
			CapturedAt: inv.Timestamp,
		})
	}

	if inv.System.Model != "" {
		snapshot.Observations = append(snapshot.Observations, model.Observation{
			DeviceID:   "system-0",
			Scope:      "hardware",
			Parameter:  "model",
			Value:      &model.Value{Text: &inv.System.Model},
			Source:     model.WMI,
			Status:     model.Observed,
			CapturedAt: inv.Timestamp,
		})
	}

	if inv.System.SerialNumber != "" {
		snapshot.Observations = append(snapshot.Observations, model.Observation{
			DeviceID:   "system-0",
			Scope:      "hardware",
			Parameter:  "serial_number",
			Value:      &model.Value{Text: &inv.System.SerialNumber},
			Source:     model.WMI,
			Status:     model.Observed,
			CapturedAt: inv.Timestamp,
		})
	}

	// Add BIOS device
	biosDevice := model.Device{
		ID:       "bios-0",
		Kind:     "bios",
		ParentID: "system-0",
	}
	snapshot.Devices = append(snapshot.Devices, biosDevice)

	// Add BIOS observations
	if inv.BIOS.Vendor != "" {
		snapshot.Observations = append(snapshot.Observations, model.Observation{
			DeviceID:   "bios-0",
			Scope:      "firmware",
			Parameter:  "vendor",
			Value:      &model.Value{Text: &inv.BIOS.Vendor},
			Source:     model.WMI,
			Status:     model.Observed,
			CapturedAt: inv.Timestamp,
		})
	}

	if inv.BIOS.Version != "" {
		snapshot.Observations = append(snapshot.Observations, model.Observation{
			DeviceID:   "bios-0",
			Scope:      "firmware",
			Parameter:  "version",
			Value:      &model.Value{Text: &inv.BIOS.Version},
			Source:     model.WMI,
			Status:     model.Observed,
			CapturedAt: inv.Timestamp,
		})
	}

	if inv.BIOS.Date != "" {
		snapshot.Observations = append(snapshot.Observations, model.Observation{
			DeviceID:   "bios-0",
			Scope:      "firmware",
			Parameter:  "date",
			Value:      &model.Value{Text: &inv.BIOS.Date},
			Source:     model.WMI,
			Status:     model.Observed,
			CapturedAt: inv.Timestamp,
		})
	}

	// Add CPU device
	if inv.CPU.Name != "" {
		cpuDevice := model.Device{
			ID:       "cpu-0",
			Kind:     "cpu",
			ParentID: "system-0",
		}
		snapshot.Devices = append(snapshot.Devices, cpuDevice)

		// Add CPU observations
		if inv.CPU.Name != "" {
			snapshot.Observations = append(snapshot.Observations, model.Observation{
				DeviceID:   "cpu-0",
				Scope:      "hardware",
				Parameter:  "name",
				Value:      &model.Value{Text: &inv.CPU.Name},
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: inv.Timestamp,
			})
		}

		if inv.CPU.Manufacturer != "" {
			snapshot.Observations = append(snapshot.Observations, model.Observation{
				DeviceID:   "cpu-0",
				Scope:      "hardware",
				Parameter:  "manufacturer",
				Value:      &model.Value{Text: &inv.CPU.Manufacturer},
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: inv.Timestamp,
			})
		}

		if inv.CPU.Architecture != "" {
			snapshot.Observations = append(snapshot.Observations, model.Observation{
				DeviceID:   "cpu-0",
				Scope:      "hardware",
				Parameter:  "architecture",
				Value:      &model.Value{Text: &inv.CPU.Architecture},
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: inv.Timestamp,
			})
		}

		if inv.CPU.Cores > 0 {
			cores := uint64(inv.CPU.Cores)
			snapshot.Observations = append(snapshot.Observations, model.Observation{
				DeviceID:   "cpu-0",
				Scope:      "hardware",
				Parameter:  "cores",
				Value:      &model.Value{Unsigned: &cores},
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: inv.Timestamp,
			})
		}

		if inv.CPU.Threads > 0 {
			threads := uint64(inv.CPU.Threads)
			snapshot.Observations = append(snapshot.Observations, model.Observation{
				DeviceID:   "cpu-0",
				Scope:      "hardware",
				Parameter:  "threads",
				Value:      &model.Value{Unsigned: &threads},
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: inv.Timestamp,
			})
		}

		if inv.CPU.MaxClockMHz > 0 {
			maxClock := uint64(inv.CPU.MaxClockMHz)
			snapshot.Observations = append(snapshot.Observations, model.Observation{
				DeviceID:   "cpu-0",
				Scope:      "hardware",
				Parameter:  "max_clock_mhz",
				Value:      &model.Value{Unsigned: &maxClock},
				Unit:       "MHz",
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: inv.Timestamp,
			})
		}
	}

	// Add memory modules
	for i, module := range inv.Memory {
		deviceID := fmt.Sprintf("memory-%d", i)
		
		memoryDevice := model.Device{
			ID:       deviceID,
			Kind:     "memory",
			ParentID: "system-0",
		}
		snapshot.Devices = append(snapshot.Devices, memoryDevice)

		// Add memory observations
		if module.DeviceLocator != "" {
			snapshot.Observations = append(snapshot.Observations, model.Observation{
				DeviceID:   deviceID,
				Scope:      "hardware",
				Parameter:  "device_locator",
				Value:      &model.Value{Text: &module.DeviceLocator},
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: inv.Timestamp,
			})
		}

		if module.BankLabel != "" {
			snapshot.Observations = append(snapshot.Observations, model.Observation{
				DeviceID:   deviceID,
				Scope:      "hardware",
				Parameter:  "bank_label",
				Value:      &model.Value{Text: &module.BankLabel},
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: inv.Timestamp,
			})
		}

		if module.Capacity > 0 {
			snapshot.Observations = append(snapshot.Observations, model.Observation{
				DeviceID:   deviceID,
				Scope:      "hardware",
				Parameter:  "capacity",
				Value:      &model.Value{Unsigned: &module.Capacity},
				Unit:       "bytes",
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: inv.Timestamp,
			})
		}

		if module.Speed > 0 {
			speed := uint64(module.Speed)
			snapshot.Observations = append(snapshot.Observations, model.Observation{
				DeviceID:   deviceID,
				Scope:      "hardware",
				Parameter:  "speed",
				Value:      &model.Value{Unsigned: &speed},
				Unit:       "MT/s",
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: inv.Timestamp,
			})
		}

		if module.Manufacturer != "" {
			snapshot.Observations = append(snapshot.Observations, model.Observation{
				DeviceID:   deviceID,
				Scope:      "hardware",
				Parameter:  "manufacturer",
				Value:      &model.Value{Text: &module.Manufacturer},
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: inv.Timestamp,
			})
		}

		if module.PartNumber != "" {
			snapshot.Observations = append(snapshot.Observations, model.Observation{
				DeviceID:   deviceID,
				Scope:      "hardware",
				Parameter:  "part_number",
				Value:      &model.Value{Text: &module.PartNumber},
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: inv.Timestamp,
			})
		}

		if module.SerialNumber != "" {
			snapshot.Observations = append(snapshot.Observations, model.Observation{
				DeviceID:   deviceID,
				Scope:      "hardware",
				Parameter:  "serial_number",
				Value:      &model.Value{Text: &module.SerialNumber},
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: inv.Timestamp,
			})
		}

		// Add form factor and memory type as string observations
		if module.FormFactor > 0 {
			formFactor := strconv.Itoa(int(module.FormFactor))
			snapshot.Observations = append(snapshot.Observations, model.Observation{
				DeviceID:   deviceID,
				Scope:      "hardware",
				Parameter:  "form_factor",
				Value:      &model.Value{Text: &formFactor},
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: inv.Timestamp,
			})
		}

		if module.MemoryType > 0 {
			memType := strconv.Itoa(int(module.MemoryType))
			snapshot.Observations = append(snapshot.Observations, model.Observation{
				DeviceID:   deviceID,
				Scope:      "hardware",
				Parameter:  "memory_type",
				Value:      &model.Value{Text: &memType},
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: inv.Timestamp,
			})
		}

		// Add SPD profile if available
		if module.SPD != nil && module.SPD.Profile != nil {
			profile := model.Profile{
				ID:       fmt.Sprintf("spd-%d", i),
				DeviceID: deviceID,
				Type:     "spd",
				Version:  "1.0",
			}

			// Add SPD observations to profile
			if module.SPD.Profile.Manufacturer != "" {
				profile.Observations = append(profile.Observations, model.Observation{
					DeviceID:   deviceID,
					Scope:      "spd",
					Parameter:  "manufacturer",
					Value:      &model.Value{Text: &module.SPD.Profile.Manufacturer},
					Source:     model.SPD,
					Status:     model.Observed,
					CapturedAt: inv.Timestamp,
				})
			}

			snapshot.Profiles = append(snapshot.Profiles, profile)
		}
	}

	return snapshot, nil
}