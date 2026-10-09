//go:build windows

package windows

import (
	"fmt"
	
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

// WMIClient provides access to Windows Management Instrumentation
type WMIClient struct {
	connection *ole.IUnknown
}

// NewWMIClient creates a new WMI client connection
func NewWMIClient() (*WMIClient, error) {
	err := ole.CoInitialize(0)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize COM: %w", err)
	}
	
	unknown, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		ole.CoUninitialize()
		return nil, fmt.Errorf("failed to create WMI locator: %w", err)
	}
	
	wmiLocator, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		unknown.Release()
		ole.CoUninitialize()
		return nil, fmt.Errorf("failed to get IDispatch interface: %w", err)
	}
	
	// Connect to local WMI service
	wmiServiceRaw, err := oleutil.CallMethod(wmiLocator, "ConnectServer")
	if err != nil {
		wmiLocator.Release()
		ole.CoUninitialize()
		return nil, fmt.Errorf("failed to connect to WMI service: %w", err)
	}
	
	wmiService := wmiServiceRaw.ToIDispatch()
	wmiLocator.Release()
	
	return &WMIClient{connection: wmiService}, nil
}

// Close releases WMI resources
func (w *WMIClient) Close() {
	if w.connection != nil {
		w.connection.Release()
		ole.CoUninitialize()
	}
}

// Query executes a WQL query and returns results
func (w *WMIClient) Query(query string) ([]map[string]interface{}, error) {
	resultRaw, err := oleutil.CallMethod(w.connection, "ExecQuery", query)
	if err != nil {
		return nil, fmt.Errorf("WMI query failed: %w", err)
	}
	
	result := resultRaw.ToIDispatch()
	defer result.Release()
	
	countVar, err := oleutil.GetProperty(result, "Count")
	if err != nil {
		return nil, fmt.Errorf("failed to get result count: %w", err)
	}
	
	count := int(countVar.Val)
	results := make([]map[string]interface{}, 0, count)
	
	enumVar, err := oleutil.CallMethod(result, "_NewEnum")
	if err != nil {
		return nil, fmt.Errorf("failed to get enumerator: %w", err)
	}
	
	enum := enumVar.ToIUnknown()
	defer enum.Release()
	
	for {
		itemRaw, err := oleutil.CallMethod(enum, "Next")
		if err != nil {
			break
		}
		
		if itemRaw.VT == ole.VT_NULL {
			break
		}
		
		item := itemRaw.ToIDispatch()
		
		// Extract properties from WMI object
		props := make(map[string]interface{})
		
		// Get Properties_ collection
		propsVar, err := oleutil.GetProperty(item, "Properties_")
		if err == nil {
			propsColl := propsVar.ToIDispatch()
			
			countVar, err := oleutil.GetProperty(propsColl, "Count")
			if err == nil {
				propCount := int(countVar.Val)
				for i := 0; i < propCount; i++ {
					propVar, err := oleutil.CallMethod(propsColl, "Item", i)
					if err != nil {
						continue
					}
					
					prop := propVar.ToIDispatch()
					
					nameVar, err := oleutil.GetProperty(prop, "Name")
					if err != nil {
						prop.Release()
						continue
					}
					
					valueVar, err := oleutil.GetProperty(prop, "Value")
					if err != nil {
						prop.Release()
						continue
					}
					
					name := nameVar.ToString()
					props[name] = valueVar.Value()
					
					prop.Release()
				}
			}
			propsColl.Release()
		}
		
		results = append(results, props)
		item.Release()
	}
	
	return results, nil
}

// SystemInfo represents basic system information
type SystemInfo struct {
	Manufacturer string
	Model        string
	SerialNumber string
	BIOSVendor   string
	BIOSVersion  string
	BIOSDate     string
}

// MemoryModule represents a physical memory module
type MemoryModule struct {
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
}

// GetSystemInfo retrieves basic system information via WMI
func GetSystemInfo() (*SystemInfo, error) {
	client, err := NewWMIClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create WMI client: %w", err)
	}
	defer client.Close()
	
	// Query computer system information
	systemResults, err := client.Query("SELECT Manufacturer, Model FROM Win32_ComputerSystem")
	if err != nil {
		return nil, fmt.Errorf("failed to query computer system: %w", err)
	}
	
	if len(systemResults) == 0 {
		return nil, fmt.Errorf("no computer system information found")
	}
	
	// Query BIOS information
	biosResults, err := client.Query("SELECT Manufacturer, SMBIOSBIOSVersion, ReleaseDate FROM Win32_BIOS")
	if err != nil {
		return nil, fmt.Errorf("failed to query BIOS: %w", err)
	}
	
	info := &SystemInfo{}
	
	// Extract system info
	if manufacturer, ok := systemResults[0]["Manufacturer"].(string); ok {
		info.Manufacturer = manufacturer
	}
	if model, ok := systemResults[0]["Model"].(string); ok {
		info.Model = model
	}
	
	// Extract BIOS info
	if len(biosResults) > 0 {
		if vendor, ok := biosResults[0]["Manufacturer"].(string); ok {
			info.BIOSVendor = vendor
		}
		if version, ok := biosResults[0]["SMBIOSBIOSVersion"].(string); ok {
			info.BIOSVersion = version
		}
		if date, ok := biosResults[0]["ReleaseDate"].(string); ok {
			info.BIOSDate = date
		}
	}
	
	return info, nil
}

// GetMemoryInfo retrieves memory module information via WMI
func GetMemoryInfo() ([]MemoryModule, error) {
	client, err := NewWMIClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create WMI client: %w", err)
	}
	defer client.Close()
	
	query := `SELECT DeviceLocator, BankLabel, Capacity, Speed, Manufacturer, 
	          PartNumber, SerialNumber, DataWidth, TotalWidth, FormFactor, 
	          MemoryType FROM Win32_PhysicalMemory`
	
	results, err := client.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query memory modules: %w", err)
	}
	
	modules := make([]MemoryModule, 0, len(results))
	
	for _, result := range results {
		module := MemoryModule{}
		
		if val, ok := result["DeviceLocator"].(string); ok {
			module.DeviceLocator = val
		}
		if val, ok := result["BankLabel"].(string); ok {
			module.BankLabel = val
		}
		if val, ok := result["Capacity"].(uint64); ok {
			module.Capacity = val
		}
		if val, ok := result["Speed"].(uint32); ok {
			module.Speed = val
		}
		if val, ok := result["Manufacturer"].(string); ok {
			module.Manufacturer = val
		}
		if val, ok := result["PartNumber"].(string); ok {
			module.PartNumber = val
		}
		if val, ok := result["SerialNumber"].(string); ok {
			module.SerialNumber = val
		}
		if val, ok := result["DataWidth"].(uint16); ok {
			module.DataWidth = val
		}
		if val, ok := result["TotalWidth"].(uint16); ok {
			module.TotalWidth = val
		}
		if val, ok := result["FormFactor"].(uint16); ok {
			module.FormFactor = val
		}
		if val, ok := result["MemoryType"].(uint16); ok {
			module.MemoryType = val
		}
		
		modules = append(modules, module)
	}
	
	return modules, nil
}