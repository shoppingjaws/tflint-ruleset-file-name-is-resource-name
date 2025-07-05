variable "new_var" {
  type = number
}

output "new_output" {
  value = "new"
}

locals {
  new_local = "new"
}

resource "aws_instance" "example" {
  ami = "ami-12345678"
}