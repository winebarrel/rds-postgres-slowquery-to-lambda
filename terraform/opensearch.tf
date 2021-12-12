resource "aws_elasticsearch_domain" "slowquery" {
  domain_name           = "slowquery"
  elasticsearch_version = "OpenSearch_1.0"

  lifecycle {
    ignore_changes = [
      vpc_options,
      advanced_options,
    ]
  }
}
