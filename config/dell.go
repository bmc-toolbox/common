package config

import (
	"encoding/json"
	"encoding/xml"
	"strings"
)

type dellVendorConfig struct {
	ConfigFormat string
	ConfigData   *dellConfig
}

type dellConfig struct {
	SystemConfiguration *dellSystemConfiguration `xml:"SystemConfiguration" json:"SystemConfiguration"`
}

type dellSystemConfiguration struct {
	XMLName    xml.Name         `xml:"SystemConfiguration"`
	Model      string           `xml:"Model,attr" json:"Model"`
	Comments   []string         `xml:"Comments>Comment,omitempty" json:"Comments,omitempty"`
	ServiceTag string           `xml:"ServiceTag,attr" json:"ServiceTag"`
	TimeStamp  string           `xml:"TimeStamp,attr" json:"TimeStamp"`
	Components []*dellComponent `xml:"Component" json:"Components"`
}

type dellComponent struct {
	XMLName    xml.Name                  `xml:"Component"`
	FQDD       string                    `xml:"FQDD,attr" json:"FQDD"`
	Attributes []*dellComponentAttribute `xml:"Attribute" json:"Attributes"`
}

type dellComponentAttribute struct {
	XMLName     xml.Name `xml:"Attribute"`
	Name        string   `xml:"Name,attr" json:"Name"`
	SetOnImport bool     `xml:"SetOnImport,omitempty" json:"SetOnImport,omitempty"`
	Comment     string   `xml:"Comment,omitempty" json:"Comment,omitempty"`
	Value       string   `xml:",chardata" json:"Value"`
}

func NewDellVendorConfigManager(configFormat string, vendorOptions map[string]string) (VendorConfigManager, error) {
	dell := &dellVendorConfig{}

	switch strings.ToLower(configFormat) {
	case "xml", "json":
		dell.ConfigFormat = strings.ToLower(configFormat)
	default:
		return nil, UnknownConfigFormatError(strings.ToLower(configFormat))
	}

	dell.ConfigData = &dellConfig{
		SystemConfiguration: &dellSystemConfiguration{},
	}

	dell.setSystemConfiguration(vendorOptions["model"], vendorOptions["servicetag"])

	return dell, nil
}

func (cm *dellVendorConfig) setSystemConfiguration(model, servicetag string) {
	cm.ConfigData.SystemConfiguration.Model = model
	cm.ConfigData.SystemConfiguration.ServiceTag = servicetag
	// TODO(jwb) Make this 'now'
	cm.ConfigData.SystemConfiguration.TimeStamp = "Tue Nov  2 21:19:16 2021"
}

// FindComponent locates an existing DellComponent if one exists in the ConfigData, if not
// it creates one and returns a pointer to that.
func (cm *dellVendorConfig) FindComponent(fqdd string) (c *dellComponent) {
	for _, c = range cm.ConfigData.SystemConfiguration.Components {
		if c.FQDD == fqdd {
			return
		}
	}

	c = &dellComponent{
		XMLName:    xml.Name{},
		FQDD:       fqdd,
		Attributes: []*dellComponentAttribute{},
	}

	cm.ConfigData.SystemConfiguration.Components = append(cm.ConfigData.SystemConfiguration.Components, c)

	return
}

// FindComponentAttribute locates an existing DellComponentAttribute if one exists in the
// ConfigData, if not it creates one and returns a pointer to that.
func (cm *dellVendorConfig) FindComponentAttribute(c *dellComponent, name string) (a *dellComponentAttribute) {
	for _, a = range c.Attributes {
		if a.Name == name {
			return
		}
	}

	a = &dellComponentAttribute{
		Name: name,
	}

	c.Attributes = append(c.Attributes, a)

	return
}

func (cm *dellVendorConfig) Raw(name, value string, menuPath []string) {
	c := cm.FindComponent(menuPath[0])
	attr := cm.FindComponentAttribute(c, name)
	attr.Value = value
}

func (cm *dellVendorConfig) Marshal() (string, error) {
	switch strings.ToLower(cm.ConfigFormat) {
	case "xml":
		x, err := xml.Marshal(cm.ConfigData.SystemConfiguration)
		if err != nil {
			return "", err
		}

		return string(x), nil
	case "json":
		x, err := json.Marshal(cm.ConfigData.SystemConfiguration)
		if err != nil {
			return "", err
		}

		return string(x), nil
	default:
		return "", UnknownConfigFormatError(strings.ToLower(cm.ConfigFormat))
	}
}

func (cm *dellVendorConfig) Unmarshal(cfgData string) (err error) {
	err = xml.Unmarshal([]byte(cfgData), cm.ConfigData)
	return
}

// Generic config options

func (cm *dellVendorConfig) EnableTPM() {
	cm.Raw("EnableTPM", "Enabled", []string{"BIOS.Setup.1-1"})
}

func (cm *dellVendorConfig) EnableSRIOV() {
	// TODO(jwb) How do we want to handle enabling this for different NICs
	cm.Raw("VirtualizationMode", "SRIOV", []string{"NIC.Slot.3-1-1"})
}
