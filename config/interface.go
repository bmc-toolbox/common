package config

import (
	"strings"

	"github.com/bmc-toolbox/common"
)

type VendorConfigManager interface {
	Raw(name, value string, menuPath []string)
	Marshal() (string, error)
	Unmarshal(cfgData string) (err error)
	StandardConfig() (biosConfig map[string]string, err error)

	BootMode(mode string) error
	BootOrder(mode string) error
	IntelSGX(mode string) error
	SecureBoot(enable bool) error
	TPM(enable bool) error
	SMT(enable bool) error
	SRIOV(enable bool) error
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
