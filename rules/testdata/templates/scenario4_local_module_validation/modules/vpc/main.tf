resource "aws_vpc" "main" {
  cidr_block = local.cidr
}
locals {
  cidr_block = "10.0.0.0/16"
}