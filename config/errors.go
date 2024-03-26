package config

import (
	"errors"
	"fmt"
)

var errUnknownConfigFormat = errors.New("unknown config format")
var errUnknownVendor = errors.New("unknown/unsupported vendor")

func UnknownConfigFormatError(format string) error {
	return fmt.Errorf("unknown config format %w : %s", errUnknownConfigFormat, format)
}

func UnknownVendorError(vendorName string) error {
	return fmt.Errorf("unknown/unsupported vendor %w : %s", errUnknownVendor, vendorName)
}
