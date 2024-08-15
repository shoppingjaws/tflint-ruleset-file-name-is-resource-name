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

## Data

Data block must be declared at `data_file_name_pattern` = `^data_.*.tf$`

:+1:
```hcl
// data_aws_instance.tf
data "aws_instance" "default"  {
  ...
}
```

:-1:
```hcl
// data.tf
data "aws_instance" "default"  {
  ...
}
```

## Resource
Anything that does not match a `Variable`,`Module`,`Provider`,`Output` or `Data` is considered a Resource.

:+1:
```hcl
// aws_instance.tf
resource "aws_instance" "default"  {
  ...
}
```

:-1:
```hcl
// ec2.tf
resource "aws_instance" "default"  {
  ...
}
```