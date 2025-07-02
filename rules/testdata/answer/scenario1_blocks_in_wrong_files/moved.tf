# Moved blocks (Terraform 1.1+)
moved {
  from = aws_instance.old_instance
  to   = aws_instance.example
}

moved {
  from = module.old_vpc
  to   = module.vpc
}