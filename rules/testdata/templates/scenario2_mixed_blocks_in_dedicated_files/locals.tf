locals {
  valid_local = "test"
}

output "should_not_be_here" {
  value = "test"
}

variable "should_not_be_here" {
  type = string
}