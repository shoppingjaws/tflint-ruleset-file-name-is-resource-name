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
- `rules/`: Individual rule implementations (currently empty after refactor)
- Rules implement `tflint.Rule` interface and analyze HCL AST
- Rules are registered in `main.go` RuleSet

**Rule Implementation Pattern:**
- Each rule analyzes Terraform configurations for naming violations
- Rules can check file names, block types, attributes, and nested blocks
- Rules emit issues with configurable severity levels
- Each rule should have corresponding `*_test.go` file

## Current State

- **Branch**: `develop/refactor-ruleset` (will merge to `main`)
- **Status**: Rules directory is empty after recent refactor - needs new rule implementations
- **Module**: Still references template path, needs updating to actual project path
- **Binary**: Currently named `tflint-ruleset-template`, needs renaming

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