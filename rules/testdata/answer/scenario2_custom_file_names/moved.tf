# Moved blocks (Terraform 1.1+)
moved {
  from = aws_instance.legacy_instance
  to   = aws_instance.example
}