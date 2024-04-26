resource "aws_route53_zone" "this" {
  name = var.host_name

  tags = var.tags
}

resource "aws_route53_record" "this" {
  zone_id = aws_route53_zone.this.zone_id
  name    = var.host_name
  type    = "NS"
  ttl     = 172800
  records = [
    aws_route53_zone.this.name_servers[0],
    aws_route53_zone.this.name_servers[1],
    aws_route53_zone.this.name_servers[2],
    aws_route53_zone.this.name_servers[3],
  ]
}

resource "aws_acm_certificate" "this" {
  domain_name               = var.host_name
  validation_method         = "DNS"
  subject_alternative_names = ["*.${var.host_name}"]

  tags = var.tags

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_route53_record" "certificate_validation" {
  for_each = {
    for dvo in aws_acm_certificate.this.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  type            = each.value.type
  zone_id         = aws_route53_zone.this.id
  records         = [each.value.record]
  ttl             = 300
}

resource "aws_acm_certificate_validation" "this" {
  certificate_arn = aws_acm_certificate.this.arn

  validation_record_fqdns = [for record in aws_route53_record.certificate_validation : record.fqdn]
}
