locals {
  host_name        = var.environment == "prod" ? var.domain_name : "${var.environment}.${var.domain_name}"
  name             = "${var.namespace}-${var.environment}-cloudfront"
  description      = var.environment == "prod" ? "CloudFront Distribution: ${var.namespace}" : "CloudFront Distribution: ${var.namespace}(${var.environment})"
  create_subdomain = var.create_subdomain && var.environment != "prod"
}

resource "aws_cloudfront_origin_access_control" "this" {
  name                              = "s3-cloudfront-oac"
  description                       = "Grant cloudfront access to s3 bucket ${aws_s3_bucket.this.id}"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

resource "aws_cloudfront_distribution" "this" {
  comment         = local.description
  enabled = true
  is_ipv6_enabled     = true
  aliases = [local.host_name]

  origin {
    domain_name = aws_s3_bucket.this.bucket_regional_domain_name
    origin_id = aws_s3_bucket.this.bucket_regional_domain_name
    origin_access_control_id = aws_cloudfront_origin_access_control.this.id
  }

  default_root_object = "index.html"

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = aws_s3_bucket.this.bucket_regional_domain_name

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 31536000
    default_ttl            = 31536000
    max_ttl                = 31536000
    compress               = true
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    acm_certificate_arn      = var.acm_certificate
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1.1_2016"
  }

  tags = merge(var.tags, {
    Bucket = "${aws_s3_bucket.this.id}"
  })
}

resource "aws_route53_zone" "env" {
  count = local.create_subdomain ? 1 : 0
  name  = local.host_name
}

resource "aws_route53_record" "env" {
  count = local.create_subdomain ? 1 : 0

  zone_id = var.zone_id
  name    = local.host_name
  type    = "NS"
  ttl     = 300
  records = [
    aws_route53_zone.env[0].name_servers[0],
    aws_route53_zone.env[0].name_servers[1],
    aws_route53_zone.env[0].name_servers[2],
    aws_route53_zone.env[0].name_servers[3]
  ]
}

## Point A record to CloudFront distribution
resource "aws_route53_record" "service_record" {
  name    = local.create_subdomain ? local.host_name : var.domain_name
  type    = "A"
  zone_id = local.create_subdomain ? aws_route53_zone.env[0].id : var.zone_id

  alias {
    name                   = aws_cloudfront_distribution.this.domain_name
    zone_id                = aws_cloudfront_distribution.this.hosted_zone_id
    evaluate_target_health = false
  }
}

