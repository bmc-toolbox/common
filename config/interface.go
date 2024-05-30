package config

import (
	"strings"

	"github.com/bmc-toolbox/common"
)

type VendorConfigManager interface {
	EnableTPM()
	EnableSRIOV()

	Raw(name, value string, menuPath []string)
	Marshal() (string, error)
	Unmarshal(cfgData string) (err error)
}

func NewVendorConfigManager(configFormat, vendorName string, vendorOptions map[string]string) (VendorConfigManager, error) {
	switch strings.ToLower(vendorName) {
	case common.VendorDell:
		return NewDellVendorConfigManager(configFormat, vendorOptions)
	case common.VendorSupermicro:
		return NewSupermicroVendorConfigManager(configFormat, vendorOptions)
	case common.VendorAsrockrack:
		return NewAsrockrackVendorConfigManager(configFormat, vendorOptions)
	default:
		return nil, UnknownVendorError(strings.ToLower(vendorName))
	}
}
