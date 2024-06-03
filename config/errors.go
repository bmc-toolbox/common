package config

import (
	"errors"
	"fmt"
)

var errUnknownConfigFormat = errors.New("unknown config format")
var errUnknownVendor = errors.New("unknown/unsupported vendor")
var errUnknownSettingType = errors.New("unknown setting type")

func UnknownConfigFormatError(format string) error {
	return fmt.Errorf("unknown config format %w : %s", errUnknownConfigFormat, format)
}

func UnknownSettingType(t string) error {
	return fmt.Errorf("unknown setting type %w : %s", errUnknownSettingType, t)
}

func UnknownVendorError(vendorName string) error {
	return fmt.Errorf("unknown/unsupported vendor %w : %s", errUnknownVendor, vendorName)
}
