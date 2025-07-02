# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a TFLint ruleset for enforcing Terraform file naming conventions. The ruleset enforces specific naming patterns:

- **Variables**: Must be in `variable.tf`
- **Modules**: Must be in `module.tf` 
- **Providers**: Must be in `provider.tf`
- **Outputs**: Must be in `output.tf`
- **Data sources**: Must follow pattern `data_*.tf` (e.g., `data_aws_instance.tf`)
- **Resources**: File name should match resource type (e.g., `aws_instance.tf` for `aws_instance` resources)
- **Locals**: Must be in `locals.tf`
- **Terraform settings**: Must be in `terraform.tf`
- **Import blocks**: Must be in `imports.tf`
- **Moved blocks**: Must be in `moved.tf`
- **Removed blocks**: Must be in `removed.tf`
- **Check blocks**: Must be in `checks.tf`
- **Ephemeral resources**: Must follow pattern `ephemeral_*.tf` (e.g., `ephemeral_aws_instance.tf`)

## Development Commands

```bash
# Build and test
make          # Default: build
make test     # Run tests with `go test ./...`
make build    # Build the plugin binary
make install  # Build and install to ~/.tflint.d/plugins

# Go commands
go test ./...     # Run all tests
go mod tidy       # Clean up dependencies
```

## Architecture

**Plugin Structure:**
- `main.go`: Plugin entry point using TFLint Plugin SDK v0.22.0
- `rules/`: Rule implementations with test coverage
  - `terraform_block_file_name_rule.go`: Main rule implementation
  - `block_manager.go`: Handles block collection and file name validation
  - `test_support.go`: Test helpers and utilities
  - `integration_test.go`: End-to-end testing with scenarios
- Rules implement `tflint.Rule` interface and analyze HCL AST
- Rules are registered in `main.go` RuleSet

**Rule Implementation Pattern:**
- Each rule analyzes Terraform configurations for naming violations
- Rules can check file names, block types, attributes, and nested blocks
- Rules emit issues with configurable severity levels
- Each rule should have corresponding `*_test.go` file

## Current State

- **Branch**: `develop/refactor-ruleset` (will merge to `main`)
- **Status**: Rules have been implemented in `rules/` directory with comprehensive test coverage
- **Module**: Updated to `github.com/shoppingjaws/tflint-ruleset-file-name-is-resource-name`
- **Binary**: Renamed to `tflint-ruleset-file-name-is-resource-name`
- **Implementation**: Full support for all Terraform block types including ephemeral resources
- **Testing**: Integration tests with scenario-based testing framework

## CI/CD

- **Build**: Runs on push/PR to main + daily cron on Ubuntu/Windows
- **Release**: Triggered on version tags, supports multi-platform builds
- **Security**: GPG signing required, SLSA attestation enabled
- **Dependencies**: Weekly Dependabot updates

## Local Testing

After implementing rules, test locally:
1. `make install` - installs plugin to `~/.tflint.d/plugins`
2. Create test Terraform files with naming violations
3. Run `tflint` to verify rules trigger correctly
```

## Memories

- TFLint plugin SDKの制限を回避するために、空ブロックの削除を追加
- Terraform 1.10で導入されたephemeralリソースのサポートを実装
- 統合テスト用のシナリオベースのテストフレームワークを構築