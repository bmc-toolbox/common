package config

import (
	"encoding/xml"
	"strings"
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
}

func NewSupermicroVendorConfigManager(configFormat string) (VendorConfigManager, error) {
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

// TODO(jwb) How do we handle the random nature of sub menus here..   we could make the user pass the explicit pointer to a menu struct, or..
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

// Generic config options

func (cm *supermicroVendorConfig) EnableTPM() {
	cm.Raw("  Security Device Support", "Enable", []string{"Trusted Computing"})
	cm.Raw("  SHA-1 PCR Bank", "Enabled", []string{"Trusted Computing"})
}

func (cm *supermicroVendorConfig) EnableSRIOV() {
	cm.Raw("SR-IOV Support", "Enabled", []string{"Advanced", "PCIe/PCI/PnP Configuration"})
}
