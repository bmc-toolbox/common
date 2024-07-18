package config

import (
	"encoding/xml"
	"strings"
)

type asrockrackVendorConfig struct {
	ConfigFormat string
	ConfigData   *asrockrackConfig
}

type asrockrackConfig struct {
	BiosCfg *asrockrackBiosCfg `xml:"BiosCfg"`
}

type asrockrackBiosCfg struct {
	XMLName xml.Name                 `xml:"BiosCfg"`
	Menus   []*asrockrackBiosCfgMenu `xml:"Menu"`
}

type asrockrackBiosCfgMenu struct {
	XMLName  xml.Name                    `xml:"Menu"`
	Name     string                      `xml:"name,attr"`
	Settings []*asrockrackBiosCfgSetting `xml:"Setting"`
	Menus    []*asrockrackBiosCfgMenu    `xml:"Menu"`
}

type asrockrackBiosCfgSetting struct {
	XMLName        xml.Name `xml:"Setting"`
	Name           string   `xml:"Name,attr"`
	Order          string   `xml:"order,attr"`
	SelectedOption string   `xml:"selectedOption,attr"`
	Type           string   `xml:"type,attr"`
}

func NewAsrockrackVendorConfigManager(configFormat string, vendorOptions map[string]string) (VendorConfigManager, error) {
	asrr := &asrockrackVendorConfig{}

	switch strings.ToLower(configFormat) {
	case "json":
		asrr.ConfigFormat = strings.ToLower(configFormat)
	default:
		return nil, UnknownConfigFormatError(strings.ToLower(configFormat))
	}

	asrr.ConfigData = &asrockrackConfig{
		BiosCfg: &asrockrackBiosCfg{},
	}

	return asrr, nil
}

// FindMenu locates an existing asrockrackBiosCfgMenu if one exists in the ConfigData, if not
// it creates one and returns a pointer to that.
func (cm *asrockrackVendorConfig) FindMenu(menuName string) (m *asrockrackBiosCfgMenu) {
	if cm.ConfigData.BiosCfg.Menus == nil {
		return
	}

	for _, m = range cm.ConfigData.BiosCfg.Menus {
		if m.Name == menuName {
			return
		}
	}

	m.Name = menuName

	cm.ConfigData.BiosCfg.Menus = append(cm.ConfigData.BiosCfg.Menus, m)

	return
}

// FindMenuSetting locates an existing asrockrackBiosCfgSetting if one exists in the
// ConfigData, if not it creates one and returns a pointer to that.
func (cm *asrockrackVendorConfig) FindMenuSetting(m *asrockrackBiosCfgMenu, name string) (s *asrockrackBiosCfgSetting) {
	for _, s = range m.Settings {
		if s.Name == name {
			return
		}
	}

	s.Name = name

	m.Settings = append(m.Settings, s)

	return
}

// TODO(jwb) How do we handle the random nature of sub menus here..   we could make the user pass the explicit pointer to a menu struct, or..
func (cm *asrockrackVendorConfig) Raw(name, value string, menuPath []string) {
}

func (cm *asrockrackVendorConfig) Marshal() (string, error) {
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

func (cm *asrockrackVendorConfig) Unmarshal(cfgData string) error {
	return xml.Unmarshal([]byte(cfgData), cm.ConfigData)
}

func (cm *asrockrackVendorConfig) StandardConfig() (biosConfig map[string]string, err error) {
	return biosConfig, err
}

// Generic config options

func (cm *asrockrackVendorConfig) BootOrder(mode string) error {
	// Unimplemented
	return nil
}

func (cm *asrockrackVendorConfig) BootMode(mode string) error {
	// Unimplemented
	return nil
}

func (cm *asrockrackVendorConfig) IntelSGX(mode string) error {
	// Unimplemented
	return nil
}

func (cm *asrockrackVendorConfig) SecureBoot(enable bool) error {
	// Unimplemented
	return nil
}

func (cm *asrockrackVendorConfig) TPM(enable bool) error {
	// Unimplemented
	return nil
}

func (cm *asrockrackVendorConfig) SMT(enable bool) error {
	// Unimplemented
	return nil
}

func (cm *asrockrackVendorConfig) SRIOV(enable bool) error {
	// Unimplemented
	return nil
}

func (cm *asrockrackVendorConfig) EnableTPM() {
	// Unimplemented
}

func (cm *asrockrackVendorConfig) EnableSRIOV() {
	// Unimplemented
}
