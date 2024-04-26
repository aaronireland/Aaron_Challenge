terraform {
  required_version = ">=1.8.0"

  required_providers {
    aws = {
      version = ">=5.35"
      source  = "hashicorp/aws"
    }
  }

  backend "s3" {
    bucket         = "sed-challenge-tfstate"
    key            = "state/sed-challenge/terraform2.tfstate"
    region         = "us-east-1"
    dynamodb_table = "app-state"
    encrypt        = true
  }
}
