resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = var.instance_type
  subnet_id     = module.vpc.public_subnets[0]
  tags          = local.common_tags
}