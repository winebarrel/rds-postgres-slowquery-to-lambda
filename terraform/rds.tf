resource "aws_db_instance" "slowquery_source" {
  identifier           = "slowquery-source"
  engine               = "postgres"
  engine_version       = "12.8"
  instance_class       = "db.t3.micro"
  allocated_storage    = 20
  parameter_group_name = aws_db_parameter_group.slowquery_source.name
  username             = "postgres"
  password             = "postgres"
  publicly_accessible  = true
  apply_immediately    = true

  enabled_cloudwatch_logs_exports = [
    "postgresql",
  ]
}

resource "aws_db_parameter_group" "slowquery_source" {
  name        = "slowquery-source"
  family      = "postgres12"
  description = "slowquery-source"

  parameter {
    name         = "log_min_duration_statement"
    value        = "100"
    apply_method = "immediate"
  }
}

resource "aws_cloudwatch_log_group" "slowquery_source" {
  name = "/aws/rds/instance/slowquery-source/postgresql"
}

output "slowquery_source_endpoint" {
  value = aws_db_instance.slowquery_source.endpoint
}
