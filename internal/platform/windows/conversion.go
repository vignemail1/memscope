package windows

import "fmt"

// convertMemoryModuleData converts WMI result data to MemoryModule with type safety
func convertMemoryModuleData(result map[string]interface{}) (MemoryModule, error) {
	module := MemoryModule{}
	
	// String fields with better error context
	if val, ok := result["DeviceLocator"].(string); ok {
		module.DeviceLocator = val
	} else if result["DeviceLocator"] != nil {
		return module, fmt.Errorf("unexpected type for DeviceLocator: %T", result["DeviceLocator"])
	}
	
	if val, ok := result["BankLabel"].(string); ok {
		module.BankLabel = val
	} else if result["BankLabel"] != nil {
		return module, fmt.Errorf("unexpected type for BankLabel: %T", result["BankLabel"])
	}
	
	// Capacity with safer type conversion
	if val := result["Capacity"]; val != nil {
		switch v := val.(type) {
		case uint64:
			module.Capacity = v
		case int64:
			if v >= 0 {
				module.Capacity = uint64(v)
			} else {
				return module, fmt.Errorf("negative capacity not allowed: %d", v)
			}
		case int32:
			if v >= 0 {
				module.Capacity = uint64(v)
			} else {
				return module, fmt.Errorf("negative capacity not allowed: %d", v)
			}
		case uint32:
			module.Capacity = uint64(v)
		default:
			return module, fmt.Errorf("unexpected type for Capacity: %T", v)
		}
	}
	
	// Speed with safer type conversion
	if val := result["Speed"]; val != nil {
		switch v := val.(type) {
		case uint32:
			module.Speed = v
		case int32:
			if v >= 0 {
				module.Speed = uint32(v)
			} else {
				return module, fmt.Errorf("negative speed not allowed: %d", v)
			}
		case uint16:
			module.Speed = uint32(v)
		case int16:
			if v >= 0 {
				module.Speed = uint32(v)
			} else {
				return module, fmt.Errorf("negative speed not allowed: %d", v)
			}
		default:
			return module, fmt.Errorf("unexpected type for Speed: %T", v)
		}
	}
	
	if val, ok := result["Manufacturer"].(string); ok {
		module.Manufacturer = val
	} else if result["Manufacturer"] != nil {
		return module, fmt.Errorf("unexpected type for Manufacturer: %T", result["Manufacturer"])
	}
	
	if val, ok := result["PartNumber"].(string); ok {
		module.PartNumber = val
	} else if result["PartNumber"] != nil {
		return module, fmt.Errorf("unexpected type for PartNumber: %T", result["PartNumber"])
	}
	
	if val, ok := result["SerialNumber"].(string); ok {
		module.SerialNumber = val
	} else if result["SerialNumber"] != nil {
		return module, fmt.Errorf("unexpected type for SerialNumber: %T", result["SerialNumber"])
	}
	
	// DataWidth with safer type conversion
	if val := result["DataWidth"]; val != nil {
		switch v := val.(type) {
		case uint16:
			module.DataWidth = v
		case int16:
			if v >= 0 {
				module.DataWidth = uint16(v)
			} else {
				return module, fmt.Errorf("negative data width not allowed: %d", v)
			}
		case uint32:
			if v <= 65535 {
				module.DataWidth = uint16(v)
			} else {
				return module, fmt.Errorf("data width too large for uint16: %d", v)
			}
		case int32:
			if v >= 0 && v <= 65535 {
				module.DataWidth = uint16(v)
			} else {
				return module, fmt.Errorf("data width out of range for uint16: %d", v)
			}
		default:
			return module, fmt.Errorf("unexpected type for DataWidth: %T", v)
		}
	}
	
	// TotalWidth with safer type conversion
	if val := result["TotalWidth"]; val != nil {
		switch v := val.(type) {
		case uint16:
			module.TotalWidth = v
		case int16:
			if v >= 0 {
				module.TotalWidth = uint16(v)
			} else {
				return module, fmt.Errorf("negative total width not allowed: %d", v)
			}
		case uint32:
			if v <= 65535 {
				module.TotalWidth = uint16(v)
			} else {
				return module, fmt.Errorf("total width too large for uint16: %d", v)
			}
		case int32:
			if v >= 0 && v <= 65535 {
				module.TotalWidth = uint16(v)
			} else {
				return module, fmt.Errorf("total width out of range for uint16: %d", v)
			}
		default:
			return module, fmt.Errorf("unexpected type for TotalWidth: %T", v)
		}
	}
	
	// FormFactor with safer type conversion
	if val := result["FormFactor"]; val != nil {
		switch v := val.(type) {
		case uint16:
			module.FormFactor = v
		case int16:
			if v >= 0 {
				module.FormFactor = uint16(v)
			} else {
				return module, fmt.Errorf("negative form factor not allowed: %d", v)
			}
		case uint32:
			if v <= 65535 {
				module.FormFactor = uint16(v)
			} else {
				return module, fmt.Errorf("form factor too large for uint16: %d", v)
			}
		case int32:
			if v >= 0 && v <= 65535 {
				module.FormFactor = uint16(v)
			} else {
				return module, fmt.Errorf("form factor out of range for uint16: %d", v)
			}
		default:
			return module, fmt.Errorf("unexpected type for FormFactor: %T", v)
		}
	}
	
	// MemoryType with safer type conversion
	if val := result["MemoryType"]; val != nil {
		switch v := val.(type) {
		case uint16:
			module.MemoryType = v
		case int16:
			if v >= 0 {
				module.MemoryType = uint16(v)
			} else {
				return module, fmt.Errorf("negative memory type not allowed: %d", v)
			}
		case uint32:
			if v <= 65535 {
				module.MemoryType = uint16(v)
			} else {
				return module, fmt.Errorf("memory type too large for uint16: %d", v)
			}
		case int32:
			if v >= 0 && v <= 65535 {
				module.MemoryType = uint16(v)
			} else {
				return module, fmt.Errorf("memory type out of range for uint16: %d", v)
			}
		default:
			return module, fmt.Errorf("unexpected type for MemoryType: %T", v)
		}
	}
	
	return module, nil
}