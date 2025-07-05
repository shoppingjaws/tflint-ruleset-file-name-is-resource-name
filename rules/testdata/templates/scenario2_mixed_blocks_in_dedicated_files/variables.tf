variable "valid_var" {
  type = string
}

resource "aws_instance" "should_not_be_here" {
  ami = "ami-12345678"
}

output "should_not_be_here" {
  value = "test"
}