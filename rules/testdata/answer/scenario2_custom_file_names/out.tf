output "instance_id" {
  description = "The ID of the EC2 instance"
  value       = aws_instance.example.id
}

output "vpc_id" {
  description = "The VPC ID"
  value       = module.vpc.vpc_id
}