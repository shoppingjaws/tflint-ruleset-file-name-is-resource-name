## Variable,Module,Provider,Output

Variable block must be declared at `variable_file_name_pattern` = `^variable.tf$`

Module block must be declared at `module_file_name_pattern` = `^module.tf$`

Provider block must be declared at `provider_file_name_pattern` = `^provider.tf$`

Output block must be declared at `output_file_name_pattern` = `^output.tf$`

:+1:
```hcl
// variable.tf
variable "name" {
  default = "test"
}
```

:+1:
```hcl
// module.tf
module "name" {
  source = "..."
}
```

:+1:
```hcl
// provider.tf
provider "aws" {
  version = "1.0"
}
```

:+1:
```hcl
// output.tf
output "name" {
  value = "name"
}
```

