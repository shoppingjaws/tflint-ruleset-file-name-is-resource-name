# Health checks
check "instance_health" {
  # Verify instance is properly configured
  assert {
    condition     = aws_instance.example.instance_type != ""
    error_message = "Instance type must be specified"
  }
  
  # Additional checks can be added here
  # assert {
  #   condition     = contains(["t2.micro", "t3.micro"], aws_instance.example.instance_type)
  #   error_message = "Instance type must be t2.micro or t3.micro"
  # }
}