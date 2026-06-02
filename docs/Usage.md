# Usage Guide

Complete reference for tfimport-cli commands, options, and examples.

**Table of Contents:**
- [Basic Commands](#basic-commands)
- [Command Reference](#command-reference)
- [Examples](#examples)
- [Global Flags](#global-flags)
- [Environment Variables](#environment-variables)
- [Terraform Cloud](#terraform-cloud)
- [Performance](#performance)
- [Troubleshooting](#troubleshooting)

---

## Basic Commands

### Analyze State

```bash
tfimport-cli analyze --state terraform.tfstate
```

Output example:
```
=== State Analysis Report ===

Summary:
  Total Resources: 42

Providers:
  - provider["registry.terraform.io/hashicorp/aws"]

Top Resource Types:
  aws_instance: 12 [supported]
  aws_security_group: 8 [supported]
  aws_subnet: 6 [supported]
  aws_iam_role: 5 [supported]
  aws_route_table: 3 [supported]

Modules:
  (root): 8 resources
  module.network: 18 resources
  module.compute: 16 resources
```

### Generate Import Blocks

```bash
# Generate all imports
tfimport-cli generate --state terraform.tfstate

# Generate for specific module
tfimport-cli generate --state terraform.tfstate --module module.network

# Preview without writing
tfimport-cli generate --state terraform.tfstate --dry-run

# Save summary report
tfimport-cli generate --state terraform.tfstate --summary report.json
```

Output file `imports.tf`:
```hcl
import {
  to = aws_vpc.main
  id = "vpc-0123456789abcdef0"
}

import {
  to = aws_subnet.public[0]
  id = "subnet-0abcdef0123456789"
}

import {
  to = module.network.aws_security_group.main
  id = "sg-0123456789abcdef0"
}
```

### Validate State

```bash
tfimport-cli validate --state terraform.tfstate
```

Validation output:
```
=== Validation Results ===

✓ State is readable
✓ State JSON is valid
✓ State version is correct (v4)
✓ Contains resources
✓ AWS resources supported: 100%

Ready for import!
```

### List Resources

```bash
# List all modules
tfimport-cli list modules --state terraform.tfstate

# List resource types
tfimport-cli list resource-types --state terraform.tfstate

# List resources in module
tfimport-cli list resources --state terraform.tfstate --module module.network

# Filter by resource type
tfimport-cli list resources --state terraform.tfstate --type aws_instance
```

---

## Command Reference

### analyze

Analyze state file and show resource inventory.

```bash
tfimport-cli analyze [OPTIONS]
```

**Options:**
- `--state PATH` — Local Terraform state file path
- `-o, --organization ORG` — Terraform Cloud organization
- `-w, --workspace WS` — Terraform Cloud workspace
- `--summary FILE` — Write JSON summary to file
- `-v, --verbose` — Verbose output
- `--help` — Show help

**Examples:**
```bash
# Local state
tfimport-cli analyze --state terraform.tfstate

# With JSON summary
tfimport-cli analyze --state terraform.tfstate --summary analysis.json

# Verbose mode
tfimport-cli analyze --state terraform.tfstate --verbose

# Terraform Cloud
tfimport-cli analyze -o my-org -w prod
```

---

### generate

Generate Terraform import blocks from state.

```bash
tfimport-cli generate [OPTIONS]
```

**Options:**
- `--state PATH` — Local Terraform state file path
- `-o, --organization ORG` — Terraform Cloud organization
- `-w, --workspace WS` — Terraform Cloud workspace
- `-m, --module MODULE` — Filter by module (root or module.name)
- `--out FILE` — Output file (default: imports.tf)
- `--summary FILE` — Write JSON summary to file
- `--dry-run` — Preview without writing to disk
- `-v, --verbose` — Verbose output
- `--help` — Show help

**Examples:**
```bash
# Generate all imports
tfimport-cli generate --state terraform.tfstate

# Generate for specific module
tfimport-cli generate --state terraform.tfstate --module module.network

# Output to custom file
tfimport-cli generate --state terraform.tfstate --out my-imports.tf

# Preview changes
tfimport-cli generate --state terraform.tfstate --dry-run

# Generate with JSON summary
tfimport-cli generate --state terraform.tfstate --summary report.json

# Terraform Cloud
tfimport-cli generate -o my-org -w prod --module module.compute

# Verbose output
tfimport-cli generate --state terraform.tfstate --verbose
```

---

### list

List resources from state.

```bash
tfimport-cli list SUBCOMMAND [OPTIONS]
```

**Subcommands:**
- `modules` — List all modules in state
- `resource-types` — List all resource types
- `resources` — List all individual resources

**Options:**
- `--state PATH` — Local Terraform state file path
- `-o, --organization ORG` — Terraform Cloud organization
- `-w, --workspace WS` — Terraform Cloud workspace
- `-m, --module MODULE` — Filter by module
- `-t, --type TYPE` — Filter by resource type
- `-v, --verbose` — Verbose output
- `--help` — Show help

**Examples:**
```bash
# List modules
tfimport-cli list modules --state terraform.tfstate

# List resource types
tfimport-cli list resource-types --state terraform.tfstate

# List all resources
tfimport-cli list resources --state terraform.tfstate

# List resources in module
tfimport-cli list resources --state terraform.tfstate --module module.network

# Filter by type
tfimport-cli list resources --state terraform.tfstate --type aws_instance

# Combined filters
tfimport-cli list resources --state terraform.tfstate --module module.network --type aws_security_group

# Verbose output
tfimport-cli list modules --state terraform.tfstate --verbose
```

---

### validate

Validate state file readiness for import.

```bash
tfimport-cli validate [OPTIONS]
```

**Options:**
- `--state PATH` — Local Terraform state file path
- `-o, --organization ORG` — Terraform Cloud organization
- `-w, --workspace WS` — Terraform Cloud workspace
- `-v, --verbose` — Verbose output
- `--help` — Show help

**Examples:**
```bash
# Local state
tfimport-cli validate --state terraform.tfstate

# Terraform Cloud
tfimport-cli validate -o my-org -w prod

# Verbose output
tfimport-cli validate --state terraform.tfstate --verbose
```

---

### version

Display version information.

```bash
tfimport-cli version
```

---

## Examples

### Scenario 1: Migrate Local State to New Repository

```bash
# 1. Analyze the state
tfimport-cli analyze --state terraform.tfstate

# 2. Validate readiness
tfimport-cli validate --state terraform.tfstate

# 3. Generate imports
tfimport-cli generate --state terraform.tfstate

# 4. Review generated imports.tf
cat imports.tf

# 5. Apply imports in target repository
cd /path/to/new/repo
terraform import -no-color < imports.tf
```

### Scenario 2: Migrate Specific Module

```bash
# 1. List modules to find the one you need
tfimport-cli list modules --state terraform.tfstate

# 2. Generate imports for that module
tfimport-cli generate --state terraform.tfstate --module module.network --out network-imports.tf

# 3. Apply in target repository
cd /path/to/new/repo
terraform import -no-color < network-imports.tf
```

### Scenario 3: Multi-Module Migration

```bash
# 1. Get list of modules
tfimport-cli list modules --state terraform.tfstate

# 2. Generate imports per module
tfimport-cli generate --state terraform.tfstate --module module.network --out network-imports.tf
tfimport-cli generate --state terraform.tfstate --module module.compute --out compute-imports.tf
tfimport-cli generate --state terraform.tfstate --module module.database --out database-imports.tf

# 3. Apply each in sequence
for file in *-imports.tf; do
  echo "Applying $file..."
  terraform import -no-color < "$file"
done
```

### Scenario 4: Terraform Cloud Integration

```bash
# 1. Export API token
export TF_API_TOKEN=<your-token>

# 2. Analyze workspace
tfimport-cli analyze -o my-org -w production

# 3. Generate imports
tfimport-cli generate -o my-org -w production --summary report.json

# 4. Review report
cat report.json | jq .
```

### Scenario 5: Large State Migration with Validation

```bash
# 1. Validate first (catches issues early)
tfimport-cli validate --state large-state.tfstate --verbose

# 2. Get counts by resource type
tfimport-cli list resource-types --state large-state.tfstate

# 3. Generate with summary
tfimport-cli generate --state large-state.tfstate --summary migration-summary.json

# 4. Review results
jq '.summary' migration-summary.json

# 5. Generate by module if needed
tfimport-cli list modules --state large-state.tfstate
tfimport-cli generate --state large-state.tfstate --module module.vpc --out vpc-imports.tf
```

### Scenario 6: Dry Run Before Actual Generation

```bash
# 1. Preview what would be generated (no files created)
tfimport-cli generate --state terraform.tfstate --dry-run

# 2. See output and verify
# (Review the import statements that would be created)

# 3. Actual generation if preview looks good
tfimport-cli generate --state terraform.tfstate
```

---

## Global Flags

Available on all commands:

| Flag | Description |
|------|-------------|
| `--state PATH` | Path to local Terraform state file |
| `-o, --organization ORG` | Terraform Cloud organization name |
| `-w, --workspace WS` | Terraform Cloud workspace name |
| `--summary FILE` | Write JSON summary to file (generate command) |
| `-v, --verbose` | Enable verbose output |
| `--help` | Show command help |
| `--version` | Show version information |

**State Source Priority:**
1. If `--state` is provided, use local state file
2. If `-o` and `-w` are provided, use Terraform Cloud
3. Error if neither is provided

---

## Environment Variables

| Variable | Description |
|----------|-------------|
| `TF_API_TOKEN` | Terraform Cloud API token (required for cloud operations) |
| `NO_COLOR` | Set to any value to disable colored output |

**Examples:**
```bash
# Disable color output
export NO_COLOR=1
tfimport-cli analyze --state terraform.tfstate

# Set API token for TFC
export TF_API_TOKEN=abcd1234...
tfimport-cli analyze -o my-org -w prod

# Temporary token
TF_API_TOKEN=xyz789 tfimport-cli generate -o my-org -w prod
```

---

## Terraform Cloud

### Authentication

1. Generate API token at `https://app.terraform.io/app/settings/tokens`
2. Export as environment variable:

```bash
export TF_API_TOKEN=<your-token>
```

### Usage

```bash
# Analyze workspace
tfimport-cli analyze -o my-organization -w my-workspace

# Generate imports
tfimport-cli generate -o my-organization -w my-workspace

# Specific module
tfimport-cli generate -o my-organization -w my-workspace -m module.network

# Validate workspace
tfimport-cli validate -o my-organization -w my-workspace

# List resources
tfimport-cli list resources -o my-organization -w my-workspace
```

### Examples

```bash
# View workspace state inventory
tfimport-cli analyze -o acme -w production

# Generate imports with summary
tfimport-cli generate -o acme -w production --summary prod-report.json

# Workspace-specific module
tfimport-cli generate -o acme -w staging -m module.database

# Verbose TFC operations
tfimport-cli validate -o acme -w production --verbose
```

---

## Performance

### State File Handling

- **Supports state files with 1000+ resources** efficiently
- Stream-capable architecture for very large states
- Standard library only — no heavy external dependencies
- Minimal memory footprint

### Optimization Tips

```bash
# For very large states, use module filtering
tfimport-cli generate --state huge-state.tfstate --module module.name

# Generate summary first to understand scope
tfimport-cli generate --state terraform.tfstate --summary report.json --dry-run

# Validate before full generation
tfimport-cli validate --state terraform.tfstate --verbose
```

---

## Troubleshooting

### Common Issues

#### 1. "State file not found"

```bash
# Verify state file path
ls -la terraform.tfstate

# Use absolute path
tfimport-cli analyze --state /absolute/path/terraform.tfstate

# Or change directory
cd /path/to/state/dir
tfimport-cli analyze --state terraform.tfstate
```

#### 2. "Invalid state file format"

```bash
# Validate state file format
terraform state list -state terraform.tfstate

# Check state version
jq '.version' terraform.tfstate
# Should be version 4

# Refresh state if needed
terraform state pull > terraform.tfstate
```

#### 3. "Terraform Cloud authentication failed"

```bash
# Verify token is set
echo $TF_API_TOKEN

# Check token validity
curl -H "Authorization: Bearer $TF_API_TOKEN" \
  https://app.terraform.io/api/v2/account/details

# Verify organization and workspace exist
tfimport-cli list modules -o your-org -w your-workspace
```

#### 4. "No resources found in state"

```bash
# Check state has resources
terraform state list -state terraform.tfstate

# Verify it's not an empty state
jq '.resources | length' terraform.tfstate
```

#### 5. "Resource type not supported" (old versions)

All AWS resources are supported in current version. If you see this:

```bash
# Upgrade to latest version
tfimport-cli version

# Download latest release
# https://github.com/njingu90/tfimport-cli/releases
```

#### 6. "Permission denied" on output file

```bash
# Check directory permissions
ls -la .

# Try different output location
tfimport-cli generate --state terraform.tfstate --out /tmp/imports.tf

# Or create with sudo
sudo tfimport-cli generate --state terraform.tfstate
```

### Debug Mode

Enable verbose output for more information:

```bash
# Analyze with verbose
tfimport-cli analyze --state terraform.tfstate --verbose

# Generate with verbose
tfimport-cli generate --state terraform.tfstate --verbose

# Validate with verbose
tfimport-cli validate --state terraform.tfstate --verbose
```

### Getting Help

```bash
# General help
tfimport-cli --help

# Command-specific help
tfimport-cli analyze --help
tfimport-cli generate --help
tfimport-cli list --help
tfimport-cli validate --help

# Version info
tfimport-cli version
```

---

## Advanced Usage

### Integration with CI/CD

```bash
# GitHub Actions example
- name: Generate Terraform Imports
  run: |
    tfimport-cli generate \
      --state terraform.tfstate \
      --summary report.json \
      --out imports.tf
  env:
    NO_COLOR: "1"

- name: Upload report
  uses: actions/upload-artifact@v2
  with:
    name: import-report
    path: report.json
```

### Scripting

```bash
#!/bin/bash

STATE_FILE="terraform.tfstate"

# Analyze
echo "Analyzing state..."
tfimport-cli analyze --state "$STATE_FILE"

# Validate
echo "Validating state..."
tfimport-cli validate --state "$STATE_FILE" || exit 1

# Generate with timestamp
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
tfimport-cli generate --state "$STATE_FILE" --out "imports_${TIMESTAMP}.tf"

echo "Complete! Generated imports_${TIMESTAMP}.tf"
```

### Batch Operations

```bash
#!/bin/bash

# Process multiple states
for state_file in *.tfstate; do
    echo "Processing: $state_file"
    
    module_name="${state_file%.tfstate}"
    tfimport-cli generate \
        --state "$state_file" \
        --out "${module_name}-imports.tf" \
        --summary "${module_name}-report.json"
done

echo "All states processed!"
```

### JSON Output Processing

```bash
# Generate summary and parse
tfimport-cli generate --state terraform.tfstate --summary report.json

# Count resources by type
jq '.resources | group_by(.type) | map({type: .[0].type, count: length})' report.json

# Extract specific resource info
jq '.resources[] | select(.type == "aws_instance")' report.json

# Generate report
jq '.summary' report.json
```

---

## See Also

- [README.md](../README.md) — Project overview
- [INSTALL.md](INSTALL.md) — Installation instructions
- [CONTRIBUTING.md](../CONTRIBUTING.md) — Contribution guidelines

## License

MIT License — See LICENSE file

## Author

njingu90

## Changelog

### v1.0.0 (2024-01-15)
- Initial release
- Support for 22 AWS resource types
- Local state file support
- Terraform Cloud workspace support
- Analyze, generate, list, validate, and version commands
- Comprehensive JSON and console reporting
