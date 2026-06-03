# Changelog

All notable changes to tfimport-cli will be documented in this file.

## [Unreleased]

### Fixed
- **aws_cloudwatch_log_subscription_filter**: Import ID now uses `|` separator (`log_group_name|name`)
- **aws_vpc_endpoint_route_table_association**: Import ID now uses `/` separator (`vpc_endpoint_id/route_table_id`)
- **Module and Type Flags**: Fixed `-module` and `-type` flags not being recognized in `generate` and `list` commands

## [1.0.0] - 2026-06-02

### Changed
- **Project Renamed**: `tfimportgen` → `tfimport-cli`
- Binary name: `tfimport-cli`
- Module path: `github.com/njingu90/tfimport-cli`

### Fixed
- **Data Sources**: Now properly skipped (mode: "data")
- **Module Resources**: Include full module path in import addresses
- **aws_vpc_dhcp_options_association**: Uses `vpc_id` only (not compound ID)
- **aws_route_table_association**: Uses `/` separator (was `:`)
- **aws_route**: Uses `_` separator (was `:`)
- **aws_security_group_rule**: Includes ALL CIDR blocks (not just first)
- **aws_network_interface_sg_attachment**: Correct order `eni_sg` (was `sg_eni`)
- **aws_iam_role_policy_attachment**: Uses `/` separator (was `:`)
- **aws_iam_user_policy_attachment**: Uses `/` separator (was `:`)
- **aws_iam_group_policy_attachment**: Uses `/` separator (was `:`)

### Features
- Generate Terraform import blocks from state files
- Support for local and Terraform Cloud state files
- Analyze state files for resource types and modules
- Validate migration readiness
- List resources by type or module
- 100+ AWS resource types supported

[Unreleased]: https://github.com/njingu90/tfimport-cli/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/njingu90/tfimport-cli/releases/tag/v1.0.0
