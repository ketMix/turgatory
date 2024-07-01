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
	BaseName    string
	Name        string         `yaml:"name"`
	Description string         `yaml:"description"`
	Type        string         `yaml:"type"`
	Professions []string       `yaml:"professions,omitempty"`
	Stats       map[string]int `yaml:"stats,omitempty"`
	Perk        string         `yaml:"perk,omitempty"`
}

// Load all equipment listed in the 'equipment/equipmentList.txt' file
func LoadEquipment() {
	// Load the equipment list
	bytes, err := FS.ReadFile("equipment/equipmentList.txt")
	if err != nil {
		panic(err)
	}

	// Parse the equipment list
	equipmentList := strings.Split(string(bytes), "\n")
	for _, name := range equipmentList {
		// Remove any whitespace
		name = strings.TrimSpace(name)

		// Lower the name for consistency
		name = strings.ToLower(name)

		// Check if the equipment is already loaded
		if _, ok := equipment[name]; ok {
			panic("Duplicate equipment listed: " + name)
		}

		// Load the equipment data from the filesystem
		path := "equipment/" + name + ".yaml"
		bytes, err := FS.ReadFile(path)
		if err != nil {
			panic("Error loading equipment yaml: " + path)

		}

		// Parse the equipment data
		var e *EquipmentAsset
		if err := yaml.Unmarshal(bytes, &e); err != nil {
			panic("Error unmarshalling equipment yaml: " + name)
		}

		e.BaseName = name
		equipment[name] = e
	}
}

func GetEquipment(name string) (*EquipmentAsset, error) {
	// Lower the name for consistency
	name = strings.ToLower(name)
	e, ok := equipment[name]
	if !ok {
		return nil, fmt.Errorf("equipment %s does not exist", name)
	}
	return e, nil
}

func GetEquipmentWithTypes(equipmentTypes []string) map[string]*EquipmentAsset {
	equipmentOfType := make(map[string]*EquipmentAsset)

	for name, e := range equipment {
		for _, t := range equipmentTypes {
			if e.Type == t {
				equipmentOfType[name] = e
			}
		}
	}
	return equipmentOfType
}
