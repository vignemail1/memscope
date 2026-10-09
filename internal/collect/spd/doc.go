// Package spd provides Serial Presence Detect (SPD) data parsing functionality
// for DDR3, DDR4, and DDR5 memory modules.
//
// SPD data contains detailed information about memory modules including:
// - JEDEC standard timing profiles
// - Intel XMP (Extreme Memory Profile) overclocking profiles  
// - AMD EXPO (Extended Profiles for Overclocking) profiles
// - Module capacity, manufacturer, and part number information
//
// The parsing implementation follows JEDEC SPD specifications:
// - DDR3: JEDEC JESD79-3F
// - DDR4: JEDEC JESD79-4C  
// - DDR5: JEDEC JESD79-5A
//
// Example usage:
//
//	spdData := readSPDFromHardware(slotIndex)
//	profile, err := spd.ParseSPDData(spdData)
//	if err != nil {
//		log.Fatalf("SPD parsing failed: %v", err)
//	}
//	
//	fmt.Printf("Memory: %s %s %dGB DDR%s\n",
//		profile.Manufacturer, profile.PartNumber,
//		profile.Capacity/(1024*1024*1024), profile.DeviceType)
//
// This package implements SPD decoding independently of hardware transport,
// requiring only the raw SPD bytes as input.
package spd
