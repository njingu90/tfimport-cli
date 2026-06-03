package cmd

import (
	"fmt"
	"os"

	"github.com/njingu90/tfimport-cli/pkg"
)

// GenerateCommand handles the generate subcommand
func GenerateCommand(gf GlobalFlags, args []string) error {
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

	if gf.Verbose {
		PrintSuccess("State loaded successfully")
	}

	// Get all resources
	resources := pkg.GetAllResources(state)

	if gf.Verbose {
		PrintInfo(fmt.Sprintf("Found %d resources in state", len(resources)))
	}

	// Apply module filter if specified
	if gf.Module != "" {
		if gf.Verbose {
			PrintInfo(fmt.Sprintf("Filtering by module: %s", gf.Module))
		}
		resources = pkg.FilterByModule(resources, gf.Module)
		if gf.Verbose {
			PrintInfo(fmt.Sprintf("After module filter: %d resources", len(resources)))
		}
	}

	// Apply resource type filter if specified
	if gf.Type != "" {
		if gf.Verbose {
			PrintInfo(fmt.Sprintf("Filtering by type: %s", gf.Type))
		}
		resources = pkg.FilterByResourceType(resources, gf.Type)
		if gf.Verbose {
			PrintInfo(fmt.Sprintf("After type filter: %d resources", len(resources)))
		}
	}

	if len(resources) == 0 {
		PrintWarning("No resources matched the filters")
		return nil
	}

	// Generate import blocks
	blocks, skipped := pkg.GenerateImportBlocksFromResources(resources)

	if gf.Verbose {
		PrintInfo(fmt.Sprintf("Generated %d import blocks, skipped %d resources", len(blocks), len(skipped)))
	}

	// Sort blocks by address
	blocks = pkg.SortImportBlocks(blocks)

	// Generate report
	report := pkg.GenerateReport(len(state.Resources), len(resources), len(blocks), len(skipped), skipped)

	// Print console report
	fmt.Print(pkg.FormatReportAsConsole(report))

	// Prepare output file path
	outputFile := gf.OutputFile
	if outputFile == "" {
		outputFile = "imports.tf"
	}

	// Write import blocks to file (unless dry-run)
	if !gf.DryRun {
		importContent := pkg.FormatImportBlocksAsFile(blocks)

		if err := os.WriteFile(outputFile, []byte(importContent), 0644); err != nil {
			PrintError(fmt.Sprintf("Failed to write import blocks: %v", err))
			return err
		}

		PrintSuccess(fmt.Sprintf("Import blocks written to %s", outputFile))
	} else {
		PrintInfo(fmt.Sprintf("Dry-run: Would write import blocks to %s", outputFile))
	}

	// Write summary file if requested
	if gf.SummaryFile != "" {
		jsonReport, err := pkg.FormatReportAsJSON(report)
		if err != nil {
			PrintError(fmt.Sprintf("Failed to format report: %v", err))
			return err
		}

		if !gf.DryRun {
			if err := os.WriteFile(gf.SummaryFile, []byte(jsonReport), 0644); err != nil {
				PrintError(fmt.Sprintf("Failed to write summary file: %v", err))
				return err
			}

			PrintSuccess(fmt.Sprintf("Summary written to %s", gf.SummaryFile))
		} else {
			PrintInfo(fmt.Sprintf("Dry-run: Would write summary to %s", gf.SummaryFile))
		}
	}

	// Print next steps
	fmt.Println()
	PrintSection("Next Steps")
	fmt.Println()

	if !gf.DryRun {
		fmt.Printf("1. Review the generated import blocks in %s\n", outputFile)
		fmt.Println("2. Add the imports to your Terraform configuration")
		fmt.Println("3. Run: terraform init")
		fmt.Println("4. Run: terraform apply -import")
		fmt.Println()
	} else {
		fmt.Println("This was a dry-run. No files were written.")
		fmt.Println("Remove --dry-run to write files to disk.")
		fmt.Println()
	}

	return nil
}
