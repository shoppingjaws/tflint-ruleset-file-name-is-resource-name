# TFLint Ruleset File name is resource name
[![Build Status](https://github.com/shoppingjaws/tflint-ruleset-file-name-is-resource-name/workflows/build/badge.svg?branch=main)](https://github.com/shoppingjaws/tflint-ruleset-file-name-is-resource-name/actions)

## Requirements

- TFLint v0.42+
- Go v1.22

## Installation

You can install the plugin with `tflint --init`. Declare a config in `.tflint.hcl` as follows:

```hcl
plugin "file-name-is-resource-name" {
  enabled = true
  source  = "github.com/shoppingjaws/tflint-ruleset-file-name-is-resource-name"
  version = "0.1.12"
}


rule "file_name_is_resource_name" {
  enabled = true
  variable_file_name_pattern= "^variable.tf$" // optional
  module_file_name_pattern = "^main.tf$" // optional
  locals_file_name_pattern= "^locals.tf$" // optional
  provider_file_name_pattern= "^provider.tf$" // optional
  output_file_name_pattern= "^output.tf$" // optional
  module_file_name_pattern= "^module.tf$" // optional
  data_file_name_pattern= "^data_.*.tf$" // optional
}
```

## Rules

|Name|Description|Severity|Enabled|Link|
| --- | --- | --- | --- | --- |
| file_name_is_resource_name | Rule for file name follows resource name |✔| [link](./docs/file_name_is_resource_name.md)|
