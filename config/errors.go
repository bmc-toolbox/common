package config

import (
	"errors"
	"fmt"
)

var errUnknownConfigFormat = errors.New("unknown config format")
var errUnknownVendor = errors.New("unknown/unsupported vendor")
var errUnknownSettingType = errors.New("unknown setting type")

var errInvalidBootModeOption = errors.New("invalid BootMode option <LEGACY|UEFI|DUAL>")
var errInvalidSGXOption = errors.New("invalid SGX option <Enabled|Disabled|Software Controlled>")

func UnknownConfigFormatError(format string) error {
	return fmt.Errorf("unknown config format %w : %s", errUnknownConfigFormat, format)
}

func UnknownSettingType(t string) error {
	return fmt.Errorf("unknown setting type %w : %s", errUnknownSettingType, t)
}

func UnknownVendorError(vendorName string) error {
	return fmt.Errorf("unknown/unsupported vendor %w : %s", errUnknownVendor, vendorName)
}

func InvalidBootModeOption(mode string) error {
	return fmt.Errorf("%w : %s", errInvalidBootModeOption, mode)
}

func InvalidSGXOption(mode string) error {
	return fmt.Errorf("%w : %s", errInvalidSGXOption, mode)
}
