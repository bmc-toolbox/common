// Package common is intended to provide a common data structure to model
// server hardware and its component attributes between libraries/tools.
package common

// Common holds attributes shared by all components
type Common struct {
	Oem         bool              `json:"oem"`
	Description string            `json:"description,omitempty"`
	Vendor      string            `json:"vendor,omitempty"`
	Model       string            `json:"model,omitempty"`
	Serial      string            `json:"serial,omitempty"`
	ProductName string            `json:"product_name,omitempty"`
	Firmware    *Firmware         `json:"firmware,omitempty"`
	Status      *Status           `json:"status,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// Device type is composed of various components
type Device struct {
	Common

	HardwareType       string               `json:"hardware_type,omitempty"`
	Chassis            string               `json:"chassis,omitempty"`
	BIOS               *BIOS                `json:"bios,omitempty"`
	BMC                *BMC                 `json:"bmc,omitempty"`
	Mainboard          *Mainboard           `json:"mainboard,omitempty"`
	CPLDs              []*CPLD              `json:"cplds"`
	TPMs               []*TPM               `json:"tpms,omitempty"`
	GPUs               []*GPU               `json:"gpus,omitempty"`
	CPUs               []*CPU               `json:"cpus,omitempty"`
	Memory             []*Memory            `json:"memory,omitempty"`
	NICs               []*NIC               `json:"nics,omitempty"`
	Drives             []*Drive             `json:"drives,omitempty"`
	StorageControllers []*StorageController `json:"storage_controller,omitempty"`
	PSUs               []*PSU               `json:"power_supplies,omitempty"`
	Enclosures         []*Enclosure         `json:"enclosures,omitempty"`
}

// NewDevice returns a pointer to an initialized Device type
func NewDevice() Device {
	return Device{
		BMC:                &BMC{NIC: &NIC{}},
		BIOS:               &BIOS{},
		Mainboard:          &Mainboard{},
		TPMs:               []*TPM{},
		CPLDs:              []*CPLD{},
		PSUs:               []*PSU{},
		NICs:               []*NIC{},
		GPUs:               []*GPU{},
		CPUs:               []*CPU{},
		Memory:             []*Memory{},
		Drives:             []*Drive{},
		StorageControllers: []*StorageController{},
		Enclosures:         []*Enclosure{},
	}
}

// Firmware struct holds firmware attributes of a device component
type Firmware struct {
	Installed  string            `json:"installed,omitempty"`
	Available  string            `json:"available,omitempty"`
	SoftwareID string            `json:"software_id,omitempty"`
	Previous   []*Firmware       `json:"previous,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// NewFirmwareObj returns a *Firmware object
func NewFirmwareObj() *Firmware {
	return &Firmware{Metadata: make(map[string]string)}
}

// Status is the health status of a component
type Status struct {
	Health         string
	State          string
	PostCode       int    `json:"post_code,omitempty"`
	PostCodeStatus string `json:"post_code_status,omitempty"`
}

// GPU component
type GPU struct {
	Common
}

// Enclosure component
type Enclosure struct {
	Common

	ID          string    `json:"id,omitempty"`
	ChassisType string    `json:"chassis_type,omitempty"`
	Firmware    *Firmware `json:"firmware,omitempty"`
}

// TPM component
type TPM struct {
	Common

	InterfaceType string `json:"interface_type,omitempty"`
}

// CPLD component
type CPLD struct {
	Common
}

// PSU component
type PSU struct {
	Common

	ID                 string `json:"id,omitempty"`
	PowerCapacityWatts int64  `json:"power_capacity_watts,omitempty"`
}

// BIOS component
type BIOS struct {
	Common

	SizeBytes     int64 `json:"size_bytes,omitempty"`
	CapacityBytes int64 `json:"capacity_bytes,omitempty" diff:"immutable"`
}

// BMC component
type BMC struct {
	Common

	ID  string `json:"id,omitempty"`
	NIC *NIC   `json:"nic,omitempty"`
}

// CPU component
type CPU struct {
	Common

	ID           string `json:"id,omitempty"`
	Slot         string `json:"slot,omitempty"`
	Architecture string `json:"architecture,omitempty"`
	ClockSpeedHz int64  `json:"clock_speeed_hz,omitempty"`
	Cores        int    `json:"cores,omitempty"`
	Threads      int    `json:"threads,omitempty"`
}

// Memory component
type Memory struct {
	Common

	ID           string `json:"id,omitempty"`
	Slot         string `json:"slot,omitempty"`
	Type         string `json:"type,omitempty"`
	SizeBytes    int64  `json:"size_bytes,omitempty"`
	FormFactor   string `json:"form_factor,omitempty"`
	PartNumber   string `json:"part_number,omitempty"`
	ClockSpeedHz int64  `json:"clock_speed_hz,omitempty"`
}

// NIC component
type NIC struct {
	Common

	ID          string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	SpeedBits   int64  `json:"speed_bits,omitempty"`
	PhysicalID  string `json:"physid,omitempty"`
	BusInfo     string `json:"bus_info,omitempty"`
	MacAddress  string `json:"macaddress,omitempty"`
}

// StorageController component
type StorageController struct {
	Common

	ID                           string `json:"id,omitempty"`
	SupportedControllerProtocols string `json:"supported_controller_protocol,omitempty"` // PCIe
	SupportedDeviceProtocols     string `json:"supported_device_protocol,omitempty"`     // Attached device protocols - SAS, SATA
	SupportedRAIDTypes           string `json:"supported_raid_types,omitempty"`
	PhysicalID                   string `json:"physid,omitempty"`
	BusInfo                      string `json:"bus_info,omitempty"`
	SpeedGbps                    int64  `json:"speed_gbps,omitempty"`
}

// Mainboard component
type Mainboard struct {
	Common

	PhysicalID string `json:"physid,omitempty"`
}

// Drive component
type Drive struct {
	Common

	ID                  string   `json:"id,omitempty"`
	OemID               string   `json:"oem_id,omitempty"`
	Type                string   `json:"drive_type,omitempty"`
	StorageController   string   `json:"storage_controller,omitempty"`
	BusInfo             string   `json:"bus_info,omitempty"`
	WWN                 string   `json:"wwn,omitempty"`
	Protocol            string   `json:"protocol,omitempty"`
	SmartStatus         string   `json:"smart_status,omitempty"`
	SmartErrors         []string `json:"smart_errors,omitempty"`
	CapacityBytes       int64    `json:"capacity_bytes,omitempty"`
	BlockSizeBytes      int64    `json:"block_size_bytes,omitempty"`
	CapableSpeedGbps    int64    `json:"capable_speed_gbps,omitempty"`
	NegotiatedSpeedGbps int64    `json:"negotiated_speed_gbps,omitempty"`
}
