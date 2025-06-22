resource "aws_iam_role" "lambda_role" {
  name = "lambda-execution-role"
}

resource "aws_s3_bucket" "wrong_place" {
  bucket = "should-be-in-s3-file"
}