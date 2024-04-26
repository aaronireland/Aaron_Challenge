locals {
  tags = merge(var.global_tags, {
    Project     = var.namespace,
    Environment = local.environment
  })
  environment = terraform.workspace == "default" ? "prod" : "${terraform.workspace}"
  name        = "${var.namespace}-${local.environment}"
}

module "dns" {
  source = "./modules/dns"

  host_name   = var.domain_name
  environment = local.environment
  region      = var.region
  tags        = local.tags
}

module "static_site" {
  source     = "./modules/static_site"
  depends_on = [module.dns]

  bucket_name_prefix = var.site_name
  acm_certificate    = module.dns.cloudfront_certificate
  namespace          = var.namespace
  environment        = local.environment
  domain_name        = var.domain_name
  zone_id            = module.dns.zone_id

  tags = local.tags
}
