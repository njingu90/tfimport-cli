package cmd

import (
	"fmt"
)

// BuildInfo contains build metadata
type BuildInfo struct {
	Version   string
	Commit    string
	BuildDate string
}

// VersionCommand handles the version subcommand
func VersionCommand(buildInfo BuildInfo) error {
	fmt.Println()
	fmt.Println("tfimport-cli - Terraform State Import Block Generator")
	fmt.Println()
	fmt.Printf("  Version:   %s\n", buildInfo.Version)
	fmt.Printf("  Commit:    %s\n", buildInfo.Commit)
	fmt.Printf("  BuildDate: %s\n", buildInfo.BuildDate)
	fmt.Println()

	return nil
}
