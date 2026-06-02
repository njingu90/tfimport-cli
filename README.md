# tfimport-cli

Generate Terraform import blocks for AWS resources from state files.

**tfimport-cli** is a read-only migration utility that helps teams migrate existing Terraform-managed infrastructure into new codebases. It analyzes state files and generates modern Terraform import blocks.

## Features

- 📖 **Read-only** — Never modifies state, applies, or infrastructure
- 🏠 **Multi-source** — Local state files and Terraform Cloud workspaces
- ✨ **Universal AWS** — All 400+ AWS resource types supported
- 🔍 **Smart filtering** — Filter by module or resource address
- 📊 **Reporting** — JSON summaries and human-readable analysis
- 🐧 **Multi-platform** — Linux, macOS, and Windows
- 🔮 **Future-proof** — New AWS services automatically supported

## AWS Resource Support

**Supports all AWS resources** (400+ resource types):
- EC2, VPC, IAM, Lambda, ECS, EKS, RDS, DynamoDB
- S3, KMS, Route53, CloudFront, ALB/NLB, API Gateway
- And 300+ more AWS services

Intelligent ID extraction works even for services not explicitly mapped.

## Installation

Download the latest [release](https://github.com/njingu90/tfimport-cli/releases):

```bash
# Linux
curl -L https://github.com/njingu90/tfimport-cli/releases/download/v1.0.0/tfimport-cli-v1.0.0-linux-amd64 -o tfimport-cli && chmod +x tfimport-cli && sudo mv tfimport-cli /usr/local/bin/

# macOS
curl -L https://github.com/njingu90/tfimport-cli/releases/download/v1.0.0/tfimport-cli-v1.0.0-darwin-arm64 -o tfimport-cli && chmod +x tfimport-cli && sudo mv tfimport-cli /usr/local/bin/

# Verify
tfimport-cli --version
```

See [INSTALL.md](docs/INSTALL.md) for detailed platform-specific instructions.

## Quick Start

```bash
# Analyze your state
tfimport-cli analyze --state terraform.tfstate

# Generate import blocks
tfimport-cli generate --state terraform.tfstate

# Preview without writing
tfimport-cli generate --state terraform.tfstate --dry-run

# Filter by module
tfimport-cli generate --state terraform.tfstate --module module.network
```

See [Usage.md](docs/Usage.md) for detailed command reference and examples.

## Terraform Cloud

Fetch state directly from Terraform Cloud:

```bash
export TF_API_TOKEN=<your-token>
tfimport-cli analyze -organization my-org -workspace prod
tfimport-cli generate -organization my-org -workspace prod
```

See [Usage.md](docs/Usage.md#terraform-cloud) for details.

## Documentation

- **[Usage Guide](docs/Usage.md)** — Command reference, advanced examples, troubleshooting
- **[Installation Guide](docs/INSTALL.md)** — Platform-specific setup instructions
- **[Contributing](CONTRIBUTING.md)** — How to contribute to the project

## License

MIT License — See [LICENSE](LICENSE) for details
