package cmd

import (
	"fmt"
	"sort"

	"github.com/njingu90/tfimport-cli/pkg"
)

// ListCommand handles the list subcommand
func ListCommand(gf GlobalFlags, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("list command requires a subcommand: modules, resources, or resource-types")
	}

	subcommand := args[0]
	subargs := args[1:]

	// Load state
	var state *pkg.TerraformState
	var err error

	if gf.StatePath != "" {
		if gf.Verbose {
			PrintInfo(fmt.Sprintf("Loading state from %s", gf.StatePath))
		}
		state, err = pkg.LoadLocalState(gf.StatePath)
	} else if gf.Organization != "" {
		if gf.Verbose {
			PrintInfo(fmt.Sprintf("Fetching state from Terraform Cloud: %s/%s", gf.Organization, gf.Workspace))
		}
		state, err = pkg.FetchStateFromTFC(gf.Organization, gf.Workspace, "")
	}

	if err != nil {
		PrintError(fmt.Sprintf("Failed to load state: %v", err))
		return err
	}

	// Validate state
	if err := pkg.ValidateState(state); err != nil {
		PrintError(fmt.Sprintf("Invalid state: %v", err))
		return err
	}

	switch subcommand {
	case "modules":
		return listModules(state, gf.Verbose)
	case "resources":
		return listResources(state, subargs, gf.Verbose)
	case "resource-types":
		return listResourceTypes(state, gf.Verbose)
	default:
		return fmt.Errorf("unknown list subcommand: %s (expected: modules, resources, or resource-types)", subcommand)
	}
}

// listModules lists all modules in the state
func listModules(state *pkg.TerraformState, verbose bool) error {
	modules := pkg.GetAllModules(state)

	if len(modules) == 0 {
		PrintWarning("No modules found in state (only root resources)")
		return nil
	}

	sort.Strings(modules)

	fmt.Println()
	PrintSection("Modules")
	fmt.Println()

	for _, module := range modules {
		fmt.Printf("  %s\n", module)
	}

	fmt.Printf("\nTotal: %d modules\n\n", len(modules))

	return nil
}

// listResourceTypes lists all resource types in the state
func listResourceTypes(state *pkg.TerraformState, verbose bool) error {
	counts := pkg.ResourceCountByType(state)

	if len(counts) == 0 {
		PrintWarning("No resources found in state")
		return nil
	}

	// Sort by type name
	var types []string
	for t := range counts {
		types = append(types, t)
	}
	sort.Strings(types)

	fmt.Println()
	PrintSection("Resource Types")
	fmt.Println()

	for _, t := range types {
		count := counts[t]
		supported := ""
		if pkg.IsSupportedAWSResource(t) {
			supported = " [supported]"
		}
		fmt.Printf("  %-40s %4d%s\n", t, count, supported)
	}

	fmt.Printf("\nTotal: %d resource types\n\n", len(types))

	return nil
}

// listResources lists resources from a specific module or with a specific type
func listResources(state *pkg.TerraformState, args []string, verbose bool) error {
	cf, err := ParseCommandFlags(args)
	if err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	var resources []pkg.Resource

	if cf.Module != "" {
		if verbose {
			PrintInfo(fmt.Sprintf("Filtering by module: %s", cf.Module))
		}
		resources = pkg.FilterByModule(pkg.GetAllResources(state), cf.Module)
	} else if cf.Type != "" {
		if verbose {
			PrintInfo(fmt.Sprintf("Filtering by type: %s", cf.Type))
		}
		resources = pkg.FilterByResourceType(pkg.GetAllResources(state), cf.Type)
	} else {
		resources = pkg.GetAllResources(state)
	}

	if len(resources) == 0 {
		PrintWarning("No resources found matching the filters")
		return nil
	}

	// Sort by address
	sort.Slice(resources, func(i, j int) bool {
		return resources[i].Address < resources[j].Address
	})

	fmt.Println()
	PrintSection("Resources")
	fmt.Println()

	// Group by type
	grouped := make(map[string][]pkg.Resource)
	for _, res := range resources {
		grouped[res.Type] = append(grouped[res.Type], res)
	}

	// Sort types
	var types []string
	for t := range grouped {
		types = append(types, t)
	}
	sort.Strings(types)

	// Print grouped
	for _, t := range types {
		fmt.Printf("  %s (%d):\n", t, len(grouped[t]))
		for _, res := range grouped[t] {
			supported := ""
			if pkg.IsSupportedAWSResource(res.Type) {
				supported = " [supported]"
			}
			fmt.Printf("    - %s%s\n", res.Address, supported)
		}
		fmt.Println()
	}

	fmt.Printf("Total: %d resources\n\n", len(resources))

	return nil
}
