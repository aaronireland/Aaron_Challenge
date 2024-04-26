resource "random_string" "random" {
  length  = 6
  special = false
  upper   = false
}

resource "aws_s3_bucket" "this" {
  bucket        = "${var.bucket_name_prefix}-${random_string.random.result}"
  force_destroy = true
}

resource "aws_s3_object" "upload_object" {
  for_each     = fileset("${path.module}/html/", "*")
  bucket       = aws_s3_bucket.this.id
  key          = each.value
  source       = "${path.module}/html/${each.value}"
  etag         = filemd5("${path.module}/html/${each.value}")
  content_type = "text/html"
}

resource "aws_s3_bucket_public_access_block" "this" {
  bucket                  = aws_s3_bucket.this.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = false
}


resource "aws_s3_bucket_policy" "this" {
  bucket = aws_s3_bucket.this.bucket
  policy = templatefile("${path.module}/templates/s3_policy.json", { 
    bucket = aws_s3_bucket.this.id,
    cloudfront = aws_cloudfront_distribution.this.arn,
  })
}
