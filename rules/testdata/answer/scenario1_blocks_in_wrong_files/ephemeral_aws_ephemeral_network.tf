ephemeral "aws_ephemeral_network" "temp_network" {
  cidr_block = "10.0.0.0/24"
  duration   = "30m"
}