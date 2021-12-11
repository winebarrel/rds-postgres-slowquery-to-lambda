terraform {
  required_version = ">= 1"

  required_providers {
    aws = {
      version = ">= 3"
      source  = "hashicorp/aws"
    }
  }
}

provider "aws" {
  region = "ap-northeast-1"
}
