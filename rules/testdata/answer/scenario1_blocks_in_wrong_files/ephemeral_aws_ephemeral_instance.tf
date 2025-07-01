ephemeral "aws_ephemeral_instance" "temp_compute" {
  instance_type = "t2.micro"
  duration      = "1h"
  
  lifecycle {
    cleanup_on_destroy = true
  }
}