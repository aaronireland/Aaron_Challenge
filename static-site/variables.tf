variable "aws_account_number" {
  type        = number
  description = "AWS account number used for deployment."
}

variable "region" {
  type        = string
  default     = "us-east-1"
  description = "Target AWS region"
}

variable "namespace" {
  type = string
  description = "Name to use as prefix for AWS resources (e.g. the project name)"
  default = "sed-challenge"
}

variable "global_tags" {
  type = map(string)
  default = {
    "ManagedBy"   = "Terraform"
  }
}

variable "site_name" {
  type    = string
  default = "static-site"
}

variable "service_docker_image" {
  type    = string
  default = "nginx:alpine"
}

variable "domain_name" {
  type    = string
}

variable "ec2_public_key" {
  type        = string
  description = "The public key data to use to generate a key pair to access ec2 instances (e.g. ssh-keygen -t rsa -b 4096 -C \"admin@example.com\")"
}

variable "custom_origin_host_header" {
  type        = string
  description = "Custom header value used by CloudFront"
  default     = "sed-challenge"
}
