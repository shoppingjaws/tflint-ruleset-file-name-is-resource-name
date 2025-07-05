# Output definitions
# This output exports the instance ID for use in other modules
output "instance_id" {
  description = "The ID of the EC2 instance"
  value       = aws_instance.example.id
  # Note: This ID can be used for monitoring setup
}