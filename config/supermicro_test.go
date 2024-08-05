package config

import (
	"testing"
)

func TestSupermicroVendorConfig_FindOrCreateMenu(t *testing.T) {
	// Create a new instance of supermicroVendorConfig
	cm := &supermicroVendorConfig{}

	// Create a slice of menus
	menus := []*supermicroBiosCfgMenu{
		{Name: "Menu1"},
		{Name: "Menu2", Menus: []*supermicroBiosCfgMenu{{Name: "SubMenu1"}}},
	}

	// Call the FindOrCreateMenu function
	menu := cm.FindOrCreateMenu(&menus, "Menu3")

	// Check if the menu was created and added to the slice
	if len(menus) != 3 {
		t.Errorf("Expected 3 menus, got: %d", len(menus))
	}

	// Check the returned menu is the same as the last menu in the slice
	if menu != menus[2] {
		t.Errorf("Expected menu: %v, got: %v", menus[2], menu)
	}

	// Call the FindOrCreateMenu function with an existing menu name
	existingMenu := cm.FindOrCreateMenu(&menus, "Menu1")

	// Check if the existing menu was found and returned
	if existingMenu != menus[0] {
		t.Errorf("Expected menu: %v, got: %v", menus[0], existingMenu)
	}

	// Call the FindOrCreateMenu function with an existing submenu name
	existingSubMenu := cm.FindOrCreateMenu(&(menus[1].Menus), "SubMenu1")

	if existingSubMenu != menus[1].Menus[0] {
		t.Errorf("Expected menu: %v, got: %v", menus[1].Menus[0], existingSubMenu)
	}
}

func TestSupermicroVendorConfig_FindOrCreateSetting(t *testing.T) {
	// Create a new instance of supermicroVendorConfig
	cm := &supermicroVendorConfig{
		ConfigData: &supermicroConfig{
			BiosCfg: &supermicroBiosCfg{
				Menus: []*supermicroBiosCfgMenu{
					{
						Name: "Menu1",
						Settings: []*supermicroBiosCfgSetting{
							{Name: "Setting1", SelectedOption: "Option1"},
						},
					},
				},
			},
		},
	}

	// Define the path and value for the setting
	path := []string{"Menu1", "Setting2"}
	value := "Option2"

	// Call the FindOrCreateSetting function
	setting := cm.FindOrCreateSetting(path, value)

	// Check if the setting was created and added to the menu
	if len(cm.ConfigData.BiosCfg.Menus[0].Settings) != 2 {
		t.Errorf("Expected 2 settings, got: %d", len(cm.ConfigData.BiosCfg.Menus[0].Settings))
	}

	// Check the returned setting is the same as the last setting in the menu
	if setting != cm.ConfigData.BiosCfg.Menus[0].Settings[1] {
		t.Errorf("Expected setting: %v, got: %v", cm.ConfigData.BiosCfg.Menus[0].Settings[1], setting)
	}

	// Call the FindOrCreateSetting function with an existing setting name
	existingSetting := cm.FindOrCreateSetting([]string{"Menu1", "Setting1"}, "Option1")

	// Check if the existing setting was found and returned
	if existingSetting != cm.ConfigData.BiosCfg.Menus[0].Settings[0] {
		t.Errorf("Expected setting: %v, got: %v", cm.ConfigData.BiosCfg.Menus[0].Settings[0], existingSetting)
	}
}
