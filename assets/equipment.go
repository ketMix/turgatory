package assets

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

var equipment = make(map[string]*EquipmentAsset)

// EquipmentAsset is an asset that represents an equipment.
// Info is stored as yaml data that represents the equipment.
type EquipmentAsset struct {
	Name        string         `yaml:"name"`
	Description string         `yaml:"description"`
	Type        string         `yaml:"type"`
	Professions []string       `yaml:"professions,omitempty"`
	Stats       map[string]int `yaml:"stats,omitempty"`
	Perk        string         `yaml:"perk,omitempty"`
	StackPath   string
}

func LoadEquipment(name string) (*EquipmentAsset, error) {
	// Lower the name for consistency
	name = strings.ToLower(name)
	// Check if the equipment is already loaded
	if equipment, ok := equipment[name]; ok {
		return equipment, nil
	}

	// Load the equipment data from the filesystem
	bytes, err := FS.ReadFile("equipment/" + name + ".yaml")
	if err != nil {
		fmt.Println("Error loading equipment yaml: ", name)
		return nil, err
	}

	// Parse the equipment data
	var e *EquipmentAsset
	if err := yaml.Unmarshal(bytes, &e); err != nil {
		fmt.Println("Error unmarshalling equipment yaml: ", name)
		return nil, err
	}

	// Set stack path
	e.StackPath = "equipment/" + strings.ToLower(e.Name)
	equipment[name] = e

	fmt.Println("Loaded equipment: ", e)
	return e, nil
}
