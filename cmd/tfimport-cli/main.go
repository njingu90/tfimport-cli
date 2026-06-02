package main

import (
	"fmt"
	"os"
	"strings"

	internalcmd "github.com/njingu90/tfimport-cli/internal/cmd"
)

// Version and build information - can be overridden at build time
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

func main() {
	// Parse global flags
	args := os.Args[1:]

	if len(args) == 0 {
		internalcmd.PrintUsage()
		os.Exit(1)
	}

	// Handle special cases
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		internalcmd.PrintUsage()
		os.Exit(0)
	}

	if args[0] == "-v" || args[0] == "--version" || args[0] == "version" {
		buildInfo := internalcmd.BuildInfo{
			Version:   Version,
			Commit:    Commit,
			BuildDate: BuildDate,
		}
		if err := internalcmd.VersionCommand(buildInfo); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Extract command
	command := args[0]
	commandArgs := args[1:]

	// Parse global flags
	gf, remainingArgs, err := internalcmd.ParseGlobalFlags(commandArgs)
	if err != nil {
		internalcmd.PrintError(fmt.Sprintf("Failed to parse flags: %v", err))
		os.Exit(1)
	}

	// Validate global flags
	if err := internalcmd.ValidateGlobalFlags(gf); err != nil {
		internalcmd.PrintError(fmt.Sprintf("Invalid flags: %v", err))
		os.Exit(1)
	}

	// Route to command handler
	var cmdErr error

	switch command {
	case "analyze":
		cmdErr = internalcmd.AnalyzeCommand(gf, remainingArgs)
	case "generate":
		cmdErr = internalcmd.GenerateCommand(gf, remainingArgs)
	case "list":
		cmdErr = internalcmd.ListCommand(gf, remainingArgs)
	case "validate":
		cmdErr = internalcmd.ValidateCommand(gf, remainingArgs)
	case "version":
		buildInfo := internalcmd.BuildInfo{
			Version:   Version,
			Commit:    Commit,
			BuildDate: BuildDate,
		}
		cmdErr = internalcmd.VersionCommand(buildInfo)
	default:
		internalcmd.PrintError(fmt.Sprintf("Unknown command: %s", command))
		internalcmd.PrintUsage()
		os.Exit(1)
	}

	if cmdErr != nil {
		if !strings.Contains(cmdErr.Error(), "Unknown") {
			internalcmd.PrintError(fmt.Sprintf("Command failed: %v", cmdErr))
		}
		os.Exit(1)
	}

	os.Exit(0)
}
