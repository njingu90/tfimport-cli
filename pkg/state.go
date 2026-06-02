package pkg

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadLocalState loads a Terraform state file from a local file path
func LoadLocalState(filepath string) (*TerraformState, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var state TerraformState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file JSON: %w", err)
	}

	return &state, nil
}

// ValidateState checks if a state is valid and has required fields
func ValidateState(state *TerraformState) error {
	if state == nil {
		return fmt.Errorf("state is nil")
	}

	if state.Version != 4 {
		return fmt.Errorf("unsupported state version: %d (expected 4)", state.Version)
	}

	if state.Resources == nil {
		return fmt.Errorf("state has no resources field")
	}

	return nil
}

// GetAllResources flattens all resources from the state, handling both root and module resources
func GetAllResources(state *TerraformState) []Resource {
	var resources []Resource

	// Process root resources
	for _, sr := range state.Resources {
		// Skip data sources - only process managed resources
		if sr.Mode == "data" {
			continue
		}

		for idx, inst := range sr.Instances {
			res := Resource{
				Type:          sr.Type,
				Name:          sr.Name,
				Module:        sr.Module,
				Address:       fmt.Sprintf("%s.%s", sr.Type, sr.Name),
				Attributes:    inst.Attributes,
				InstanceIndex: inst.IndexKey,
			}

			// Build address with module prefix if exists
			var baseAddr string
			if sr.Module != "" {
				baseAddr = fmt.Sprintf("%s.%s.%s", sr.Module, sr.Type, sr.Name)
			} else {
				baseAddr = fmt.Sprintf("%s.%s", sr.Type, sr.Name)
			}
			res.Address = baseAddr

			// Handle count/for_each indices
			if inst.IndexKey != nil {
				if numIdx, ok := inst.IndexKey.(float64); ok {
					res.Address = fmt.Sprintf("%s[%d]", baseAddr, int(numIdx))
				} else if strIdx, ok := inst.IndexKey.(string); ok {
					res.Address = fmt.Sprintf("%s[%q]", baseAddr, strIdx)
				}
			}

			// If only one instance and no index key, don't add index
			if len(sr.Instances) == 1 && inst.IndexKey == nil {
				res.Address = baseAddr
			} else if inst.IndexKey == nil {
				res.Address = fmt.Sprintf("%s[%d]", baseAddr, idx)
			}

			resources = append(resources, res)
		}
	}

	return resources
}

// GetResourcesByModule returns all resources under a specific module
func GetResourcesByModule(state *TerraformState, moduleName string) []Resource {
	var resources []Resource

	for _, sr := range state.Resources {
		// Skip data sources - only process managed resources
		if sr.Mode == "data" {
			continue
		}

		// Check if this resource belongs to the module
		if sr.Module == moduleName {
			for idx, inst := range sr.Instances {
				res := Resource{
					Type:          sr.Type,
					Name:          sr.Name,
					Module:        moduleName,
					Address:       fmt.Sprintf("%s.%s.%s", moduleName, sr.Type, sr.Name),
					Attributes:    inst.Attributes,
					InstanceIndex: inst.IndexKey,
				}

				// Handle count/for_each indices
				if inst.IndexKey != nil {
					if numIdx, ok := inst.IndexKey.(float64); ok {
						res.Address = fmt.Sprintf("%s.%s.%s[%d]", moduleName, sr.Type, sr.Name, int(numIdx))
					} else if strIdx, ok := inst.IndexKey.(string); ok {
						res.Address = fmt.Sprintf("%s.%s.%s[%q]", moduleName, sr.Type, sr.Name, strIdx)
					}
				}

				// If only one instance and no index key, don't add index
				if len(sr.Instances) == 1 && inst.IndexKey == nil {
					res.Address = fmt.Sprintf("%s.%s.%s", moduleName, sr.Type, sr.Name)
				} else if inst.IndexKey == nil {
					res.Address = fmt.Sprintf("%s.%s.%s[%d]", moduleName, sr.Type, sr.Name, idx)
				}

				resources = append(resources, res)
			}
		}
	}

	return resources
}

// GetAllModules returns all unique module paths in the state
func GetAllModules(state *TerraformState) []string {
	moduleMap := make(map[string]bool)

	for _, sr := range state.Resources {
		// Skip data sources
		if sr.Mode == "data" {
			continue
		}
		if sr.Module != "" {
			moduleMap[sr.Module] = true
		}
	}

	modules := make([]string, 0, len(moduleMap))
	for m := range moduleMap {
		modules = append(modules, m)
	}

	return modules
}

// GetResourceCount returns the count of managed resources (excluding data sources)
func GetResourceCount(state *TerraformState) int {
	count := 0
	for _, sr := range state.Resources {
		if sr.Mode != "data" {
			count++
		}
	}
	return count
}

// ResourceCountByType returns a map of resource type to count
func ResourceCountByType(state *TerraformState) map[string]int {
	counts := make(map[string]int)

	for _, sr := range state.Resources {
		// Skip data sources
		if sr.Mode == "data" {
			continue
		}
		counts[sr.Type]++
	}

	return counts
}

// GetProvidersInUse returns all unique provider names in the state
func GetProvidersInUse(state *TerraformState) []string {
	providerMap := make(map[string]bool)

	for _, sr := range state.Resources {
		if sr.Provider != "" {
			providerMap[sr.Provider] = true
		}
	}

	providers := make([]string, 0, len(providerMap))
	for p := range providerMap {
		providers = append(providers, p)
	}

	return providers
}
