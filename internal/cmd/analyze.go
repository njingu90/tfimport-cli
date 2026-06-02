package cmd

import (
	"fmt"
	"os"

	"github.com/njingu90/tfimport-cli/pkg"
)

// AnalyzeCommand handles the analyze subcommand
func AnalyzeCommand(gf GlobalFlags, args []string) error {
	// Load state
	var state *pkg.TerraformState
	var source string
	var err error

	if gf.StatePath != "" {
		if gf.Verbose {
			PrintInfo(fmt.Sprintf("Loading state from %s", gf.StatePath))
		}
		state, err = pkg.LoadLocalState(gf.StatePath)
		source = "local state"
	} else if gf.Organization != "" {
		if gf.Verbose {
			PrintInfo(fmt.Sprintf("Fetching state from Terraform Cloud: %s/%s", gf.Organization, gf.Workspace))
		}
		state, err = pkg.FetchStateFromTFC(gf.Organization, gf.Workspace, "")
		source = fmt.Sprintf("Terraform Cloud (%s/%s)", gf.Organization, gf.Workspace)
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
		PrintSuccess(fmt.Sprintf("Loaded state from %s", source))
	}

	// Generate analysis report
	report := pkg.GenerateAnalysisReport(state)

	// Print console report
	fmt.Print(pkg.FormatAnalysisReportAsConsole(report))

	// Write summary file if requested
	if gf.SummaryFile != "" {
		jsonReport, err := pkg.FormatReportAsJSON(report)
		if err != nil {
			PrintError(fmt.Sprintf("Failed to format report: %v", err))
			return err
		}

		if err := os.WriteFile(gf.SummaryFile, []byte(jsonReport), 0644); err != nil {
			PrintError(fmt.Sprintf("Failed to write summary file: %v", err))
			return err
		}

		PrintSuccess(fmt.Sprintf("Summary written to %s", gf.SummaryFile))
	}

	return nil
}
