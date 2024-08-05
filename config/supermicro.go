package config

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"

	"golang.org/x/net/html/charset"
)

const (
	// enabledValue and disabledValue are utilized for bios setting value normalization
	enabledValue  = "Enabled"
	disabledValue = "Disabled"
)

type supermicroVendorConfig struct {
	ConfigFormat string
	ConfigData   *supermicroConfig
}

type supermicroConfig struct {
	BiosCfg *supermicroBiosCfg `xml:"BiosCfg"`
}

type supermicroBiosCfg struct {
	XMLName xml.Name                 `xml:"BiosCfg"`
	Menus   []*supermicroBiosCfgMenu `xml:"Menu"`
}

type supermicroBiosCfgMenu struct {
	XMLName  xml.Name                    `xml:"Menu"`
	Name     string                      `xml:"name,attr"`
	Settings []*supermicroBiosCfgSetting `xml:"Setting"`
	Menus    []*supermicroBiosCfgMenu    `xml:"Menu"`
}

type supermicroBiosCfgSetting struct {
	XMLName        xml.Name `xml:"Setting"`
	Name           string   `xml:"name,attr"`
	Order          string   `xml:"order,attr"`
	SelectedOption string   `xml:"selectedOption,attr"`
	Type           string   `xml:"type,attr"`
	CheckedStatus  string   `xml:"checkedStatus,attr"`
	NumericValue   string   `xml:"numericValue,attr"`
}

func NewSupermicroVendorConfigManager(configFormat string, vendorOptions map[string]string) (VendorConfigManager, error) {
	supermicro := &supermicroVendorConfig{}

	switch strings.ToLower(configFormat) {
	case "xml":
		supermicro.ConfigFormat = strings.ToLower(configFormat)
	default:
		return nil, UnknownConfigFormatError(strings.ToLower(configFormat))
	}

	supermicro.ConfigData = &supermicroConfig{
		BiosCfg: &supermicroBiosCfg{},
	}

	return supermicro, nil
}

// Function to find or create a setting by path
func (cm *supermicroVendorConfig) FindOrCreateSetting(path []string, value string) *supermicroBiosCfgSetting {
	biosCfg := cm.ConfigData.BiosCfg

	var currentMenus = &biosCfg.Menus

	for i, part := range path {
		if i == len(path)-1 {
			// Last part, create or find the setting
			for j := range *currentMenus {
				for k := range (*currentMenus)[j].Settings {
					if (*currentMenus)[j].Settings[k].Name == part {
						return (*currentMenus)[j].Settings[k]
					}
				}
			}

			// If no setting found in any menu, create a new setting in the first menu
			newSetting := supermicroBiosCfgSetting{Name: part, SelectedOption: ""}
			(*currentMenus)[0].Settings = append((*currentMenus)[0].Settings, &newSetting)

			return (*currentMenus)[0].Settings[len((*currentMenus)[0].Settings)-1]
		} else {
			// Intermediate part, find or create the menu
			_ = cm.FindOrCreateMenu(currentMenus, part)
		}
	}

	return nil
}

// Function to find or create a menu by name
func (cm *supermicroVendorConfig) FindOrCreateMenu(menus *[]*supermicroBiosCfgMenu, name string) *supermicroBiosCfgMenu {
	for i := range *menus {
		if (*menus)[i].Name == name {
			return (*menus)[i]
		}
	}

	newMenu := &supermicroBiosCfgMenu{Name: name}
	*menus = append(*menus, newMenu)

	return (*menus)[len(*menus)-1]
}

func (cm *supermicroVendorConfig) Raw(name, value string, menuPath []string) {
	menuPath = append(menuPath, name)

	_ = cm.FindOrCreateSetting(menuPath, value)
}

func (cm *supermicroVendorConfig) Marshal() (string, error) {
	switch strings.ToLower(cm.ConfigFormat) {
	case "xml":
		x, err := xml.Marshal(cm.ConfigData)
		if err != nil {
			return "", err
		}

		return string(x), nil
	default:
		return "", UnknownConfigFormatError(strings.ToLower(cm.ConfigFormat))
	}
}

func (cm *supermicroVendorConfig) Unmarshal(cfgData string) (err error) {
	// the xml exported by sum is ISO-8859-1 encoded
	decoder := xml.NewDecoder(bytes.NewReader([]byte(cfgData)))
	// convert characters from non-UTF-8 to UTF-8
	decoder.CharsetReader = charset.NewReaderLabel

	return decoder.Decode(cm.ConfigData.BiosCfg)
}

func (cm *supermicroVendorConfig) StandardConfig() (biosConfig map[string]string, err error) {
	biosConfig = make(map[string]string)

	for _, menu := range cm.ConfigData.BiosCfg.Menus {
		for _, s := range menu.Settings {
			switch s.Name {
			// We want to drop this list of settings
			case "NewSetupPassword", "NewSysPassword", "OldSetupPassword", "OldSysPassword":
			// All others get normalized
			default:
				var k, v string
				k, v, err = normalizeSetting(s)

				if err != nil {
					return
				}

				biosConfig[k] = v
			}
		}
	}

	return biosConfig, err
}

func normalizeSetting(s *supermicroBiosCfgSetting) (k, v string, err error) {
	switch s.Type {
	case "CheckBox":
		k = normalizeName(s.Name)
		v = normalizeValue(k, s.CheckedStatus)
	case "Option":
		k = normalizeName(s.Name)
		v = normalizeValue(k, s.SelectedOption)
	case "Password":
		k = normalizeName(s.Name)
		v = ""
	case "Numeric":
		k = normalizeName(s.Name)
		v = normalizeValue(k, s.NumericValue)
	default:
		err = UnknownSettingType(s.Type)
		return
	}

	return
}

func normalizeName(k string) string {
	switch k {
	case "CpuMinSevAsid":
		return "amd_sev"
	case "BootMode", "Boot mode select":
		return "boot_mode"
	case "IntelTxt":
		return "intel_txt"
	case "Software Guard Extensions (SGX)":
		return "intel_sgx"
	case "SecureBoot", "Secure Boot":
		return "secure_boot"
	case "Hyper-Threading", "Hyper-Threading [ALL]", "LogicalProc":
		return "smt"
	case "SriovGlobalEnable":
		return "sr_iov"
	case "TpmSecurity", "Security Device Support":
		return "tpm"
	default:
		// When we don't normalize the key prepend "raw:"
		return "raw:" + k
	}
}

func normalizeBootMode(v string) string {
	switch strings.ToLower(v) {
	case "legacy":
		return "BIOS"
	default:
		return strings.ToUpper(v)
	}
}

func normalizeValue(k, v string) string {
	if k == "boot_mode" {
		return normalizeBootMode(v)
	}

	switch strings.ToLower(v) {
	case "disable":
		return disabledValue
	case "disabled":
		return disabledValue
	case "enable":
		return enabledValue
	case "enabled":
		return enabledValue
	case "off":
		return disabledValue
	case "on":
		return enabledValue
	default:
		return v
	}
}

// Generic config options

func (cm *supermicroVendorConfig) BootMode(mode string) error {
	switch strings.ToUpper(mode) {
	case "LEGACY", "UEFI", "DUAL":
		cm.Raw("Boot mode select", strings.ToUpper(mode), []string{"Boot"})
	default:
		return InvalidBootModeOption(strings.ToUpper(mode))
	}

	return nil
}

func (cm *supermicroVendorConfig) BootOrder(mode string) error {
	// In a supermicro config there are 8 total legacy boot options and 9 UEFI boot options
	// Since we primarily care about the first two boot options we explicitly define them
	// and rely on the for loop to populate the remainder as Disabled.
	switch strings.ToUpper(mode) {
	case "LEGACY":
		cm.Raw("Legacy Boot Option #1", "Hard Disk", []string{"Boot"})
		cm.Raw("Legacy Boot Option #2", "Network", []string{"Boot"})

		for i := 3; i < 8; i++ {
			cm.Raw("Legacy Boot Option #"+fmt.Sprint(i), "Disabled", []string{"Boot"})
		}
	case "UEFI":
		cm.Raw("UEFI Boot Option #1", "UEFI Hard Disk", []string{"Boot"})
		cm.Raw("UEFI Boot Option #2", "UEFI Network", []string{"Boot"})

		for i := 3; i < 9; i++ {
			cm.Raw("UEFI Boot Option #"+fmt.Sprint(i), "Disabled", []string{"Boot"})
		}
	case "DUAL":
		// TODO(jwb) Is this just both sets?
	default:
		return InvalidBootModeOption(strings.ToUpper(mode))
	}

	return nil
}

func (cm *supermicroVendorConfig) IntelSGX(mode string) error {
	switch mode {
	case "Disabled", "Enabled", "Software Controlled":
		// TODO(jwb) Path needs to be determined.
		cm.Raw("Software Guard Extensions (SGX)", mode, []string{"Advanced", "PCIe/PCI/PnP Configuration"})
	default:
		return InvalidSGXOption(mode)
	}

	return nil
}

func (cm *supermicroVendorConfig) SecureBoot(enable bool) error {
	if !enable {
		cm.Raw("Secure Boot", "Disabled", []string{"SMC Secure Boot Configuration"})
	} else {
		cm.Raw("Secure Boot", "Enabled", []string{"SMC Secure Boot Configuration"})
	}

	return nil
}

func (cm *supermicroVendorConfig) TPM(enable bool) error {
	if enable {
		// Note, this is actually 'Enable' not 'Enabled' like everything else.
		cm.Raw(" Security Device Support", "Enable", []string{"Trusted Computing"})
		cm.Raw(" SHA-1 PCR Bank", "Enabled", []string{"Trusted Computing"})
	} else {
		// Note, this is actually 'Disable' not 'Disabled' like everything else.
		cm.Raw(" Security Device Support", "Disable", []string{"Trusted Computing"})
		cm.Raw(" SHA-1 PCR Bank", "Disabled", []string{"Trusted Computing"})
	}

	return nil
}

func (cm *supermicroVendorConfig) SMT(enable bool) error {
	if enable {
		cm.Raw("Hyper-Threading", "Enabled", []string{"Advanced", "CPU Configuration"})
	} else {
		cm.Raw("Hyper-Threading", "Disabled", []string{"Advanced", "CPU Configuration"})
	}

	return nil
}

func (cm *supermicroVendorConfig) SRIOV(enable bool) error {
	// TODO(jwb) Need to figure out how we do this on platforms that support it...
	return nil
}
