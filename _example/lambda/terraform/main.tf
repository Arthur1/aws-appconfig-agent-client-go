terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "ap-northeast-1"
}

locals {
  application_name               = "demo-app-lambda"
  environment_name               = "demo-env"
  configuration1_name            = "demo-conf1"
  configuration2_name            = "demo-conf2"
  configuration_featureflag_name = "demo-featureflag"
}
