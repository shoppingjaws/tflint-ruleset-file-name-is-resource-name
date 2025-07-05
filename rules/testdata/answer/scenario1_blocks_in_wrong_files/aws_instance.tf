resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = var.instance_type
  tags          = local.common_tags
}