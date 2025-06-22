variable "new_var" {
  type = number
}

resource "aws_instance" "example" {
  ami = "ami-12345678"
}