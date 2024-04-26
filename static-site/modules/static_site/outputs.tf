output "s3_bucket_website_endpoint" {
  # value = aws_s3_bucket_website_configuration.this.website_endpoint
  value = aws_s3_bucket.this.bucket_regional_domain_name
}

output "s3_bucket_name" {
  value = aws_s3_bucket.this.id
}

output "cloudfront_domain_name" {
  value = aws_cloudfront_distribution.this.domain_name
}

output "domain_name" {
  value = aws_route53_record.service_record.name
}
