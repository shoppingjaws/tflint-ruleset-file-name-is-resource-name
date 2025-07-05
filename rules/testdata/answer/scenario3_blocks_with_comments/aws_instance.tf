# Main EC2 instance resource
# This creates the primary application server
resource "aws_instance" "example" {
  ami           = "ami-12345678" # Amazon Linux 2
  instance_type = var.instance_type
  tags          = local.common_tags
  
  # TODO: Add user_data script for initialization
  # user_data = file("init.sh")
}