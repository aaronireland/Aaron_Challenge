output "zone_id" {
  value = aws_route53_zone.this.zone_id
}

output "acm_certificate" {
  value = aws_acm_certificate_validation.this.certificate_arn
}

output "cloudfront_certificate" {
  value = aws_acm_certificate_validation.this.certificate_arn
}
