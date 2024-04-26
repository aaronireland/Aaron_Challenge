variable "namespace" {
  type = string
}

variable "environment" {
  type = string
}

variable "domain_name" {
  type = string
}

variable create_subdomain {
  type = bool
  default = false
}

variable "zone_id" {
  type = string
}

variable "tags" {
  type = map(string)
}

variable "acm_certificate" {
  type        = string
  description = "ACM Certificate ARN for CloudFront in us-east-1"
}
variable "bucket_name_prefix" {
  type = string
}

variable "index_document" {
  type = string
  default = "index.html"
}
