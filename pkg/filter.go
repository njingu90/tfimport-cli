package pkg

import (
	"fmt"
	"strings"
)

// FilterByModule returns resources from a specific module
func FilterByModule(resources []Resource, moduleName string) []Resource {
	var filtered []Resource

	for _, res := range resources {
		if res.Module == moduleName {
			filtered = append(filtered, res)
		}
	}

	return filtered
}

// FilterByResourceType returns resources matching a specific type
func FilterByResourceType(resources []Resource, resourceType string) []Resource {
	var filtered []Resource

	for _, res := range resources {
		if res.Type == resourceType {
			filtered = append(filtered, res)
		}
	}

	return filtered
}

// FilterByResourceAddress returns a single resource matching an exact address
// Supports addresses like:
// - aws_vpc.main
// - module.network.aws_vpc.main
// - aws_vpc.main[0]
// - module.network.aws_vpc.main["key"]
func FilterByResourceAddress(resources []Resource, address string) *Resource {
	// Normalize address format
	for _, res := range resources {
		if res.Address == address {
			return &res
		}
	}

	return nil
}

// FilterByPrefix returns all resources whose address starts with a prefix
// Useful for finding all resources in a module hierarchy
func FilterByPrefix(resources []Resource, prefix string) []Resource {
	var filtered []Resource

	for _, res := range resources {
		if strings.HasPrefix(res.Address, prefix) {
			filtered = append(filtered, res)
		}
	}

	return filtered
}

// FilterByResourceTypesSupported returns only resources with supported types
func FilterByResourceTypesSupported(resources []Resource, supportedTypes map[string]bool) ([]Resource, []Resource) {
	var supported, unsupported []Resource

	for _, res := range resources {
		if supportedTypes[res.Type] {
			supported = append(supported, res)
		} else {
			unsupported = append(unsupported, res)
			res.SkipReason = fmt.Sprintf("resource type %s not supported in v1.0", res.Type)
			unsupported[len(unsupported)-1] = res
		}
	}

	return supported, unsupported
}

// GetTopResourceTypes returns the N most common resource types
func GetTopResourceTypes(resources []Resource, limit int) []struct {
	Type  string
	Count int
} {
	counts := make(map[string]int)

	for _, res := range resources {
		counts[res.Type]++
	}

	// Convert to sorted slice
	type TypeCount struct {
		Type  string
		Count int
	}

	var types []TypeCount
	for t, c := range counts {
		types = append(types, TypeCount{Type: t, Count: c})
	}

	// Simple sort (bubble sort for small lists)
	for i := 0; i < len(types); i++ {
		for j := i + 1; j < len(types); j++ {
			if types[j].Count > types[i].Count {
				types[i], types[j] = types[j], types[i]
			}
		}
	}

	var result []struct {
		Type  string
		Count int
	}

	for i := 0; i < limit && i < len(types); i++ {
		result = append(result, struct {
			Type  string
			Count int
		}{Type: types[i].Type, Count: types[i].Count})
	}

	return result
}

// ExtractModulePrefix extracts the module path from a resource address
// e.g., "module.network.module.compute.aws_vpc.main" -> "module.network.module.compute"
func ExtractModulePrefix(address string) string {
	parts := strings.Split(address, ".")

	// Find the last occurrence of "aws_" or other provider prefixes
	for i := len(parts) - 1; i >= 0; i-- {
		if strings.HasPrefix(parts[i], "aws_") || strings.Contains(parts[i], "_") {
			// Check if this is a resource type
			if i >= 1 && (parts[i-1] == "module" || (i > 1 && parts[i-2] == "module")) {
				// Return everything before the resource type
				return strings.Join(parts[:i], ".")
			}
		}
	}

	return ""
}
