data "aws_iam_role" "postgresql_slowquery" {
  name = "postgresql-slowquery-role-gwzngy7v"
}

resource "aws_ecr_repository" "postgresql_slowquery" {
  name = "postgresql-slowquery"
}

output "postgresql_slowquery_repository_url" {
  value = aws_ecr_repository.postgresql_slowquery.repository_url
}

resource "aws_lambda_function" "postgresql_slowquery" {
  function_name = "postgresql-slowquery"
  role          = data.aws_iam_role.postgresql_slowquery.arn
  package_type  = "Image"
  image_uri     = "${aws_ecr_repository.postgresql_slowquery.repository_url}:latest"

  environment {
    variables = {
      DD_API_KEY = var.dd_api_key
      DD_APP_KEY = var.dd_app_key
    }
  }
}

resource "aws_lambda_permission" "invoke_lambda" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.postgresql_slowquery.function_name
  principal     = "logs.ap-northeast-1.amazonaws.com"
  source_arn    = "${aws_cloudwatch_log_group.slowquery_source.arn}:*"
}

resource "aws_cloudwatch_log_subscription_filter" "postgresql_slowquery" {
  name            = "LambdaStream_${aws_lambda_function.postgresql_slowquery.function_name}"
  distribution    = "ByLogStream"
  log_group_name  = aws_cloudwatch_log_group.slowquery_source.name
  filter_pattern  = ""
  destination_arn = aws_lambda_function.postgresql_slowquery.arn
}

resource "aws_cloudwatch_log_group" "postgresql_slowquery" {
  name = "/aws/lambda/postgresql-slowquery"
}
