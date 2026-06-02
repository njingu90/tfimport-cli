package cmd

import (
	"fmt"
	"strings"

	"github.com/njingu90/tfimport-cli/pkg"
)

// ValidateCommand handles the validate subcommand
func ValidateCommand(gf GlobalFlags, args []string) error {
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

	// Perform validation checks
	checks := performValidationChecks(state)

	// Print results
	fmt.Println()
	PrintSection("Validation Results")
	fmt.Println()

	passCount := 0
	failCount := 0

	for _, check := range checks {
		if check.passed {
			PrintSuccess(check.message)
			passCount++
		} else {
			PrintWarning(check.message)
			failCount++
		}
	}

	fmt.Println()

	// Summary
	if failCount > 0 {
		PrintWarning(fmt.Sprintf("Validation failed: %d issue(s) found", failCount))
		fmt.Println()
		return fmt.Errorf("validation failed")
	}

	PrintSuccess(fmt.Sprintf("All validation checks passed (%d checks)", passCount))
	fmt.Println()

	// Readiness assessment
	fmt.Println()
	PrintSection("Migration Readiness")
	fmt.Println()

	readinessChecks := performReadinessChecks(state)

	for _, check := range readinessChecks {
		if check.ready {
			PrintSuccess(check.message)
		} else {
			PrintWarning(check.message)
		}
	}

	fmt.Println()

	return nil
}

// ValidationCheck represents the result of a validation check
type ValidationCheck struct {
	message string
	passed  bool
}

// ReadinessCheck represents the result of a readiness assessment
type ReadinessCheck struct {
	message string
	ready   bool
}

// performValidationChecks runs all validation checks on a state
func performValidationChecks(state *pkg.TerraformState) []ValidationCheck {
	var checks []ValidationCheck

	// Check: State is not nil
	if state == nil {
		checks = append(checks, ValidationCheck{
			message: "State is readable",
			passed:  false,
		})
		return checks
	}

	checks = append(checks, ValidationCheck{
		message: "State is readable",
		passed:  true,
	})

	// Check: JSON is valid (already parsed)
	checks = append(checks, ValidationCheck{
		message: "State JSON is valid",
		passed:  true,
	})

	// Check: Version is correct
	versionOk := state.Version == 4
	checks = append(checks, ValidationCheck{
		message: fmt.Sprintf("State version is correct (v%d)", state.Version),
		passed:  versionOk,
	})

	// Check: Resources are present
	hasResources := len(state.Resources) > 0
	checks = append(checks, ValidationCheck{
		message: fmt.Sprintf("State has resources (%d resources)", len(state.Resources)),
		passed:  hasResources,
	})

	// Check: Supported resources
	supportedCount := 0
	unsupportedCount := 0
	for _, sr := range state.Resources {
		if pkg.IsSupportedAWSResource(sr.Type) {
			supportedCount++
		} else if strings.HasPrefix(sr.Type, "aws_") {
			unsupportedCount++
		}
	}

	checks = append(checks, ValidationCheck{
		message: fmt.Sprintf("Supported AWS resources: %d (unsupported: %d)", supportedCount, unsupportedCount),
		passed:  supportedCount > 0,
	})

	// Check: No obviously broken resources (with missing IDs)
	missingIDCount := 0
	for _, sr := range state.Resources {
		if pkg.IsSupportedAWSResource(sr.Type) {
			res := pkg.Resource{
				Type:       sr.Type,
				Name:       sr.Name,
				Attributes: sr.Instances[0].Attributes,
			}
			if err := pkg.ValidateResourceForImport(res); err != nil {
				missingIDCount++
			}
		}
	}

	checks = append(checks, ValidationCheck{
		message: fmt.Sprintf("No resources with missing import IDs (%d OK)", supportedCount-missingIDCount),
		passed:  missingIDCount == 0,
	})

	return checks
}

// performReadinessChecks assesses migration readiness
func performReadinessChecks(state *pkg.TerraformState) []ReadinessCheck {
	var checks []ReadinessCheck

	// Assessment: AWS provider only
	providers := pkg.GetProvidersInUse(state)
	providersOk := true
	for _, p := range providers {
		if !strings.HasPrefix(p, "aws") {
			providersOk = false
			break
		}
	}

	checks = append(checks, ReadinessCheck{
		message: fmt.Sprintf("Using AWS provider only (%s)", strings.Join(providers, ", ")),
		ready:   providersOk,
	})

	// Assessment: Resource count
	resourceCount := len(state.Resources)

	topResources := pkg.GetTopResourceTypes(pkg.GetAllResources(state), 3)
	topList := []string{}
	for _, t := range topResources {
		topList = append(topList, fmt.Sprintf("%s (%d)", t.Type, t.Count))
	}

	checks = append(checks, ReadinessCheck{
		message: fmt.Sprintf("State has %d resources (top types: %s)", resourceCount, strings.Join(topList, ", ")),
		ready:   resourceCount > 0,
	})

	// Assessment: Module structure
	modules := pkg.GetAllModules(state)
	if len(modules) > 0 {
		checks = append(checks, ReadinessCheck{
			message: fmt.Sprintf("Organized in %d modules", len(modules)),
			ready:   true,
		})
	} else {
		checks = append(checks, ReadinessCheck{
			message: "Resources are in root module only",
			ready:   true,
		})
	}

	// Assessment: Supported resource coverage
	supportedCount := 0
	totalCount := 0
	coverage := 0
	for _, sr := range state.Resources {
		if strings.HasPrefix(sr.Type, "aws_") {
			totalCount++
			if pkg.IsSupportedAWSResource(sr.Type) {
				supportedCount++
			}
		}
	}

	if totalCount > 0 {
		coverage = (supportedCount * 100) / totalCount
		ready := coverage >= 80 // At least 80% covered
		checks = append(checks, ReadinessCheck{
			message: fmt.Sprintf("AWS resource coverage is %d%% (%d/%d)", coverage, supportedCount, totalCount),
			ready:   ready,
		})
	}

	// Recommendation
	fmt.Println()
	fmt.Println("Recommendation:")
	if providersOk && resourceCount > 0 {
		if coverage >= 80 {
			fmt.Println("  ✓ State is ready for import. Run: tfimportgen generate -state terraform.tfstate")
		} else {
			fmt.Println("  ⚠ Some resources are not yet supported. You may need to manually import them.")
			fmt.Println("    Recommendation: Generate what you can, then manually handle unsupported resources.")
		}
	} else if !providersOk {
		fmt.Println("  ✗ State uses non-AWS providers. tfimportgen currently supports AWS only.")
	} else {
		fmt.Println("  ✗ State is empty or invalid.")
	}

	return checks
}
