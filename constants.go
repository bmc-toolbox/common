package common

import "strings"

const (
	VendorDell                  = "dell"
	VendorMicron                = "micron"
	VendorAsrockrack            = "asrockrack"
	VendorSupermicro            = "supermicro"
	VendorHPE                   = "hp"
	VendorQuanta                = "quanta"
	VendorGigabyte              = "gigabyte"
	VendorIntel                 = "intel"
	VendorLSI                   = "lsi"
	VendorHGST                  = "hgst"
	VendorPacket                = "packet"
	VendorMellanox              = "mellanox"
	VendorToshiba               = "toshiba"
	VendorAmericanMegatrends    = "ami"
	VendorBroadcom              = "broadcom"
	VendorInfineon              = "infineon"
	SystemManufacturerUndefined = "To Be Filled By O.E.M."

	// Generic component slugs
	// Slugs are set on Device types to identify the type of component
	SlugBackplaneExpander     = "Backplane-Expander"
	SlugChassis               = "Chassis"
	SlugTPM                   = "TPM"
	SlugGPU                   = "GPU"
	SlugCPU                   = "CPU"
	SlugPhysicalMem           = "PhysicalMemory"
	SlugStorageController     = "StorageController"
	SlugStorageControllers    = "StorageControllers"
	SlugBMC                   = "BMC"
	SlugBIOS                  = "BIOS"
	SlugDrive                 = "Drive"
	SlugDrives                = "Drives"
	SlugDriveTypePCIeNVMEeSSD = "NVMe-PCIe-SSD"
	SlugDriveTypeSATASSD      = "Sata-SSD"
	SlugDriveTypeSATAHDD      = "Sata-HDD"
	SlugNIC                   = "NIC"
	SlugNICs                  = "NICs"
	SlugPSU                   = "Power-Supply"
	SlugPSUs                  = "Power-Supplies"
	SlugCPLD                  = "CPLD"
	SlugEnclosure             = "Enclosure"
	SlugMainboard             = "Mainboard"
	SlugUnknown               = "unknown"

	// Smart status
	SmartStatusOK      = "ok"
	SmartStatusFailed  = "failed"
	SmartStatusUnknown = "unknown"
)

// downcases and returns a normalized vendor name from the given string
func FormatVendorName(v string) string {
	switch v {
	case "ASRockRack":
		return VendorAsrockrack
	case "Dell Inc.":
		return VendorDell
	case "HP", "HPE":
		return VendorHPE
	case "Supermicro":
		return VendorSupermicro
	case "Quanta Cloud Technology Inc.":
		return VendorQuanta
	case "GIGABYTE":
		return VendorGigabyte
	case "Intel Corporation":
		return VendorIntel
	case "Packet":
		return VendorPacket
	default:
		return v
	}
}

// Return the product vendor name, given a product name/model string
func VendorFromString(s string) string {
	s = strings.ToLower(s)

	switch {
	case strings.Contains(s, "dell"):
		return VendorDell
	case strings.Contains(s, "lsi3008-it"):
		return VendorLSI
	case strings.Contains(s, "hgst "):
		return VendorHGST
	case strings.Contains(s, "intel "):
		return VendorIntel
	case strings.Contains(s, "micron_"), strings.HasPrefix(s, "mtfd"):
		return VendorMicron
	case strings.Contains(s, "toshiba"):
		return VendorToshiba
	case strings.Contains(s, "connectx4lx"):
		return VendorMellanox
	case strings.Contains(s, "infineon"):
		return VendorInfineon
	default:
		return ""
	}
}
