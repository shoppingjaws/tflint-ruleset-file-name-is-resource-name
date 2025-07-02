# TFLint Ruleset: File Name is Resource Name

[![Build Status](https://github.com/shoppingjaws/tflint-ruleset-file-name-is-resource-name/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/shoppingjaws/tflint-ruleset-file-name-is-resource-name/actions)
[![GitHub release](https://img.shields.io/github/release/shoppingjaws/tflint-ruleset-file-name-is-resource-name.svg)](https://github.com/shoppingjaws/tflint-ruleset-file-name-is-resource-name/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

A TFLint ruleset for enforcing Terraform file naming conventions. This plugin ensures your Terraform code follows consistent file naming patterns, making your codebase more organized and maintainable.

## 📋 Overview

This ruleset enforces specific naming patterns for different Terraform block types:

| Block Type | Required File Name | Example |
|------------|-------------------|----------|
| Variables | `variable.tf` | All `variable` blocks must be in this file |
| Modules | `module.tf` | All `module` blocks must be in this file |
| Providers | `provider.tf` | All `provider` blocks must be in this file |
| Outputs | `output.tf` | All `output` blocks must be in this file |
| Locals | `locals.tf` | All `locals` blocks must be in this file |
| Terraform settings | `terraform.tf` | `terraform` block must be in this file |
| Import blocks | `imports.tf` | All `import` blocks must be in this file |
| Moved blocks | `moved.tf` | All `moved` blocks must be in this file |
| Removed blocks | `removed.tf` | All `removed` blocks must be in this file |
| Check blocks | `checks.tf` | All `check` blocks must be in this file |
| Data sources | `data_*.tf` | `data_aws_instance.tf` for `data "aws_instance"` blocks |
| Resources | `<resource_type>.tf` | `aws_instance.tf` for `resource "aws_instance"` blocks |
| Ephemeral resources | `ephemeral_*.tf` | `ephemeral_aws_instance.tf` for `ephemeral "aws_instance"` blocks |

## 🚀 Features

- **Comprehensive Coverage**: Supports all Terraform block types including the latest ephemeral resources (Terraform 1.10+)
- **Clear Error Messages**: Provides specific guidance on where each block should be placed
- **Flexible Configuration**: Each rule can be individually enabled/disabled
- **Performance**: Efficient AST-based analysis for fast linting
- **Well-Tested**: Extensive test coverage with scenario-based integration tests

## 📦 Installation

### Prerequisites

- TFLint v0.42+
- Go v1.24+ (for building from source)

### Using TFLint's Plugin System

You can install the plugin with `tflint --init`. Add this to your `.tflint.hcl`:

```hcl
plugin "file-name-is-resource-name" {
  enabled = true
  version = "0.1.0"
  source  = "github.com/shoppingjaws/tflint-ruleset-file-name-is-resource-name"
}
```

### Building from Source

```bash
# Clone the repository
git clone https://github.com/shoppingjaws/tflint-ruleset-file-name-is-resource-name.git
cd tflint-ruleset-file-name-is-resource-name

# Build the plugin
make build

# Install to ~/.tflint.d/plugins
make install
```

## 📝 Configuration

Create a `.tflint.hcl` file in your Terraform project:

```hcl
plugin "file-name-is-resource-name" {
  enabled = true
}

# You can disable specific rules if needed
rule "terraform_block_file_name" {
  enabled = false
}
```

## 🎯 Usage

After installation, run TFLint in your Terraform project:

```bash
tflint
```

Example output when violations are found:

```
1 issue(s) found:

Error: Resource block for "aws_instance" must be in file "aws_instance.tf" (terraform_block_file_name)
  on main.tf line 1:
   1: resource "aws_instance" "example" {
```

## 🛠️ Development

### Project Structure

```
.
├── main.go                 # Plugin entry point
├── rules/                  # Rule implementations
│   ├── terraform_block_file_name_rule.go
│   ├── block_manager.go    # Core logic for block validation
│   └── *_test.go          # Comprehensive test coverage
└── testdata/              # Test scenarios
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestScenario ./rules
```

### Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [TFLint](https://github.com/terraform-linters/tflint) for providing an excellent plugin framework
- The Terraform community for feedback and contributions

## 📚 Resources

- [TFLint Plugin Development Guide](https://github.com/terraform-linters/tflint/blob/master/docs/developer-guide/plugins.md)
- [Terraform Documentation](https://www.terraform.io/docs)
- [Project Issues](https://github.com/shoppingjaws/tflint-ruleset-file-name-is-resource-name/issues)
