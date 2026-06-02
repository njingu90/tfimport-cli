package cmd

import (
	"flag"
	"fmt"
	"os"
)

// GlobalFlags represents global command-line flags
type GlobalFlags struct {
	StatePath    string
	Organization string
	Workspace    string
	OutputFile   string
	SummaryFile  string
	DryRun       bool
	Verbose      bool
}

// CommandFlags represents command-specific flags
type CommandFlags struct {
	Module string
	Type   string
}

// ParseGlobalFlags parses global flags from arguments
func ParseGlobalFlags(args []string) (GlobalFlags, []string, error) {
	fs := flag.NewFlagSet("global", flag.ContinueOnError)

	var gf GlobalFlags
	fs.StringVar(&gf.StatePath, "state", "", "Path to local Terraform state file")
	fs.StringVar(&gf.Organization, "organization", "", "Terraform Cloud organization")
	fs.StringVar(&gf.Workspace, "workspace", "", "Terraform Cloud workspace")
	fs.StringVar(&gf.OutputFile, "out", "", "Output file for generated import blocks")
	fs.StringVar(&gf.SummaryFile, "summary", "", "Output file for JSON summary report")
	fs.BoolVar(&gf.DryRun, "dry-run", false, "Perform analysis without writing files")
	fs.BoolVar(&gf.Verbose, "verbose", false, "Enable verbose output")

	// Parse flags, stopping at first non-flag argument
	err := fs.Parse(args)
	if err != nil {
		return GlobalFlags{}, nil, err
	}

	return gf, fs.Args(), nil
}

// ParseCommandFlags parses command-specific flags
func ParseCommandFlags(args []string) (CommandFlags, error) {
	fs := flag.NewFlagSet("command", flag.ContinueOnError)

	var cf CommandFlags
	fs.StringVar(&cf.Module, "module", "", "Filter by module (e.g., module.network)")
	fs.StringVar(&cf.Type, "type", "", "Filter by resource type")

	err := fs.Parse(args)
	if err != nil {
		return CommandFlags{}, err
	}

	return cf, nil
}

// ValidateGlobalFlags checks that required global flags are set
func ValidateGlobalFlags(gf GlobalFlags) error {
	// Must have either local state or TFC connection
	if gf.StatePath == "" && gf.Organization == "" {
		return fmt.Errorf("either --state or --organization must be specified")
	}

	// If TFC, both org and workspace are required
	if (gf.Organization != "" && gf.Workspace == "") || (gf.Organization == "" && gf.Workspace != "") {
		return fmt.Errorf("both --organization and --workspace must be specified together")
	}

	return nil
}

// GetStateSource determines the state source from global flags
func GetStateSource(gf GlobalFlags) (*StateSource, error) {
	if gf.StatePath != "" {
		return &StateSource{
			Type:      "local",
			LocalPath: gf.StatePath,
		}, nil
	}

	if gf.Organization != "" && gf.Workspace != "" {
		return &StateSource{
			Type:         "tfc",
			TFCOrg:       gf.Organization,
			TFCWorkspace: gf.Workspace,
		}, nil
	}

	return nil, fmt.Errorf("no state source specified")
}

// StateSource represents a Terraform state source
type StateSource struct {
	Type         string
	LocalPath    string
	TFCOrg       string
	TFCWorkspace string
}

// PrintUsage prints usage information
func PrintUsage() {
	usage := `Usage: tfimport-cli <command> [options]

Commands:
  analyze     Analyze and report on state file
  generate    Generate Terraform import blocks
  list        List resources from state
  validate    Validate state and migration readiness
  version     Show version information

Global Options:
  -state string           Path to local Terraform state file
  -organization string    Terraform Cloud organization
  -workspace string       Terraform Cloud workspace
  -out string            Output file for generated imports (default: imports.tf)
  -summary string        Output file for JSON summary
  -dry-run               Perform analysis without writing files
  -verbose               Enable verbose output

Examples:
  tfimport-cli analyze -state terraform.tfstate
  tfimport-cli generate -state terraform.tfstate -module module.network
  tfimport-cli list modules -state terraform.tfstate
  tfimport-cli validate -organization my-org -workspace prod
`
	fmt.Fprint(os.Stderr, usage)
}
