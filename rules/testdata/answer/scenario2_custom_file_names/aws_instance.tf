resource "aws_instance" "example" {
  ami           = "ami-0c02fb55956c7d316"
  instance_type = var.instance_type

  vpc_security_group_ids = [aws_security_group.example.id]
  subnet_id              = module.vpc.public_subnets[0]

  tags = local.common_tags
}