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
}

func NewVendorConfigManager(configFormat, vendorName string) (VendorConfigManager, error) {
	switch strings.ToLower(vendorName) {
	case common.VendorDell:
		return NewDellVendorConfigManager(configFormat)
	case common.VendorSupermicro:
		return NewSupermicroVendorConfigManager(configFormat)
	case common.VendorAsrockrack:
		return NewAsrockrackVendorConfigManager(configFormat)
	default:
		return nil, UnknownVendorError(strings.ToLower(vendorName))
	}
}
