output "valid_output" {
  value = "test"
}

variable "should_not_be_here" {
  type = string
}

locals {
  should_not_be_here = "test"
}