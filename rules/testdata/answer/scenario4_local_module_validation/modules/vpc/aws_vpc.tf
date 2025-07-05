resource "aws_vpc" "main" {
  cidr_block = local.cidr_block
}