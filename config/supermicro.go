package config

import (
	"encoding/xml"
	"strings"
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
	Name           string   `xml:"Name,attr"`
	Order          string   `xml:"order,attr"`
	SelectedOption string   `xml:"selectedOption,attr"`
	Type           string   `xml:"type,attr"`
	CheckedStatus  string   `xml:"checkedStatus,attr"`
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

// FindMenu locates an existing SupermicroBiosCfgMenu if one exists in the ConfigData, if not
// it creates one and returns a pointer to that.
func (cm *supermicroVendorConfig) FindMenu(menuName string, menuRoot *supermicroBiosCfgMenu) (m *supermicroBiosCfgMenu) {
	// root is cm.ConfigData.BiosCfg.Menus
	for _, m = range menuRoot.Menus {
		if m.Name == menuName {
			return
		}
	}

	m.Name = menuName

	menuRoot.Menus = append(menuRoot.Menus, m)

	return
}

// FindMenuSetting locates an existing SupermicroBiosCfgSetting if one exists in the
// ConfigData, if not it creates one and returns a pointer to that.
func (cm *supermicroVendorConfig) FindMenuSetting(m *supermicroBiosCfgMenu, name string) (s *supermicroBiosCfgSetting) {
	for _, s = range m.Settings {
		if s.Name == name {
			return
		}
	}

	s.Name = name

	m.Settings = append(m.Settings, s)

	return
}

func (cm *supermicroVendorConfig) Raw(name, value string, menuPath []string) {
	menus := make([]*supermicroBiosCfgMenu, 0, len(menuPath))

	for i, name := range menuPath {
		var m *supermicroBiosCfgMenu

		if i == 0 {
			m = cm.FindMenu(name, cm.ConfigData.BiosCfg.Menus[0])
		} else {
			m = cm.FindMenu(name, menus[i-1])
		}

		menus = append(menus, m)
	}
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
	err = xml.Unmarshal([]byte(cfgData), cm.ConfigData)
	return
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

	return
}

func normalizeSetting(s *supermicroBiosCfgSetting) (k, v string, err error) {
	switch s.Type {
	case "CheckBox":
		k = normalizeName(s.Name)
		v = normalizeValue(k, s.CheckedStatus)
		return
	case "Option":
		k = normalizeName(s.Name)
		v = normalizeValue(k, s.SelectedOption)
		return
	default:
		err = UnknownSettingType(s.Type)
		return
	}
}

func normalizeName(name string) string {
	switch name {
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
		// When we don't normalize the value append "raw:" to the value
		return "raw:" + name
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

func (cm *supermicroVendorConfig) EnableTPM() {
	cm.Raw("  Security Device Support", "Enable", []string{"Trusted Computing"})
	cm.Raw("  SHA-1 PCR Bank", "Enabled", []string{"Trusted Computing"})
}

func (cm *supermicroVendorConfig) EnableSRIOV() {
	cm.Raw("SR-IOV Support", "Enabled", []string{"Advanced", "PCIe/PCI/PnP Configuration"})
}
