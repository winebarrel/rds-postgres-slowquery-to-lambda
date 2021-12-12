resource "aws_elasticsearch_domain" "es_slowquery" {
  domain_name           = "es-slowquery"
  elasticsearch_version = "7.10"

  lifecycle {
    ignore_changes = [
      vpc_options,
      advanced_options,
    ]
  }
}
