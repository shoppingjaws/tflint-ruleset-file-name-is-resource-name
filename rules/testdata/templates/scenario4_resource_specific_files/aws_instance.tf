resource "aws_instance" "correct" {
  ami = "ami-12345678"
}

resource "aws_s3_bucket" "wrong_type" {
  bucket = "wrong-resource-type"
}

variable "should_not_be_here" {
  type = string
}