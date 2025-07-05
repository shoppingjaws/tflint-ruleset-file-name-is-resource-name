# Local module references - these should trigger file naming rules
module "vpc" {
  source = "./modules/vpc"
  
  environment = "production"
  cidr_block  = "10.0.0.0/16"
}
