
resource "aws_s3_bucket" "this" {
  bucket = var.bucket
}

resource "aws_s3_bucket_public_access_block" "this" {
  bucket                  = aws_s3_bucket.this.bucket
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

locals {
  this_object_path              = "archive/bootstrap.zip"
  this_object_base64sha256_path = "${local.this_object_path}.base64sha256"
  this_function_dir             = "${path.module}/src/${var.function.name}/"
  this_function_code = join("", [
    for file in fileset(local.this_function_dir, "**")
    :filebase64sha256("${local.this_function_dir}/${file}") if endswith(file,".go") || endswith(file,"go.mod") || endswith(file,"go.sum") 
    
  ])
}

resource "null_resource" "this" {
  triggers = {
    code_diff = "${local.this_function_code}"
  }
  provisioner "local-exec" {
    command = "cd ${local.this_function_dir} && GOOS=linux GOARCH=amd64 go build -o ./build/bootstrap"
  }
  provisioner "local-exec" {
    command = "cd ${local.this_function_dir} && zip -j ./archive/bootstrap.zip ./build/bootstrap"
  }
  provisioner "local-exec" {
    command = "cd ${local.this_function_dir} && aws s3 cp ./${local.this_object_path} s3://${aws_s3_bucket.this.bucket}/${local.this_object_path} --checksum-algorithm SHA256"
  }
  provisioner "local-exec" {
    command = "cd ${local.this_function_dir} && openssl dgst -sha256 -binary ./${local.this_object_path} | openssl enc -base64 | tr -d \"\n\" > ./${local.this_object_base64sha256_path}"
  }
  provisioner "local-exec" {
    command = "cd ${local.this_function_dir} && aws s3 cp ./${local.this_object_base64sha256_path} s3://${aws_s3_bucket.this.bucket}/${local.this_object_base64sha256_path} --content-type \"text/plain\" --checksum-algorithm SHA256"
  }
}

data "aws_s3_object" "zip" {
  bucket     = aws_s3_bucket.this.bucket
  key        = local.this_object_path
  depends_on = [null_resource.this]
}

data "aws_s3_object" "zip_base64sha256" {
  bucket     = aws_s3_bucket.this.bucket
  key        = local.this_object_base64sha256_path
  depends_on = [null_resource.this]
}

resource "aws_lambda_function" "this" {
  function_name    = var.function.name
  handler          = "main.lambda_handler"
  runtime          = "provided.al2023"
  s3_bucket        = aws_s3_bucket.this.bucket
  s3_key           = data.aws_s3_object.zip.key
  source_code_hash = data.aws_s3_object.zip_base64sha256.body
  role             = aws_iam_role.this.arn
  timeout          = 30
  publish          = true
  environment {
    variables = var.env
  }
}

resource "aws_lambda_function_url" "this" {
  function_name      = aws_lambda_function.this.function_name
  authorization_type = "NONE"
}

resource "aws_lambda_function_event_invoke_config" "this" {
  function_name                = aws_lambda_function.this.function_name
  maximum_event_age_in_seconds = 60 * 10 # minutes
  maximum_retry_attempts       = 0
}

resource "aws_iam_role" "this" {
  name = "${var.function.name}-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = ["sts:AssumeRole"]
        Principal = {
          Service = ["lambda.amazonaws.com"]
        }
      }
    ]
  })
}
