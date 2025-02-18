# terraform {
#   required_providers {
#     aws = {
#       source  = "hashicorp/aws"
#       version = "5.85.0"
#     }
#   }

#   backend "s3" {
#     bucket  = "localstack-demo-terraform-state"
#     region  = "eu-west-2"
#     key     = "localstack-demo.tfstate"
#     encrypt = true
#   }
# }

provider "aws" {
  region                      = "eu-west-2"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
  endpoints {
    apigateway      = "http://localhost.localstack.cloud:4566"
    apigatewayv2    = "http://localhost.localstack.cloud:4566"
    cloudformation  = "http://localhost.localstack.cloud:4566"
    cloudwatch      = "http://localhost.localstack.cloud:4566"
    cognitoidp      = "http://localhost.localstack.cloud:4566"
    cognitoidentity = "http://localhost.localstack.cloud:4566"
    dynamodb        = "http://localhost.localstack.cloud:4566"
    ec2             = "http://localhost.localstack.cloud:4566"
    es              = "http://localhost.localstack.cloud:4566"
    elasticache     = "http://localhost.localstack.cloud:4566"
    firehose        = "http://localhost.localstack.cloud:4566"
    iam             = "http://localhost.localstack.cloud:4566"
    kinesis         = "http://localhost.localstack.cloud:4566"
    lambda          = "http://localhost.localstack.cloud:4566"
    rds             = "http://localhost.localstack.cloud:4566"
    redshift        = "http://localhost.localstack.cloud:4566"
    route53         = "http://localhost.localstack.cloud:4566"
    s3              = "http://s3.localhost.localstack.cloud:4566"
    secretsmanager  = "http://localhost.localstack.cloud:4566"
    ses             = "http://localhost.localstack.cloud:4566"
    sns             = "http://localhost.localstack.cloud:4566"
    sqs             = "http://localhost.localstack.cloud:4566"
    ssm             = "http://localhost.localstack.cloud:4566"
    stepfunctions   = "http://localhost.localstack.cloud:4566"
    sts             = "http://localhost.localstack.cloud:4566"
  }
}

module "localstack-demo-lambda" {
  source = "../../modules/lambda"
  bucket = "user-lambda"
  function = {
    name = "user-lambda"
  }
  env = {
    DYNAMODB_TABLE_NAME = aws_dynamodb_table.this.name
    DYNAMODB_ENDPOINT   = "http://localhost.localstack.cloud:4566"
    USER_SQS_URL        = aws_sqs_queue.test.url
    SQS_ENDPOINT        = "http://localhost.localstack.cloud:4566"
    DB_NAME             = "ih_authenticator"
    DB_USER             = "root"
    DB_PASSWORD         = "password"
    DB_READER_HOST      = "db.lh.local"
    DB_WRITER_HOST      = "db.lh.local"
    DB_PORT             = "3306"
  }
}

resource "aws_dynamodb_table" "this" {
  name         = "users"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"
  attribute {
    name = "id"
    type = "S"
  }
}

resource "aws_sqs_queue" "test" {
  name                       = "test"
  visibility_timeout_seconds = 10
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.test_dead_letter.arn
    maxReceiveCount     = 2
  })
}

resource "aws_sqs_queue" "test_dead_letter" {
  name = "test-dead-letter"
}

resource "aws_lambda_event_source_mapping" "event_source_mapping" {
  event_source_arn                   = aws_sqs_queue.test.arn
  enabled                            = true
  function_name                      = module.sqs_lambda.function_arn
  batch_size                         = 2
  maximum_batching_window_in_seconds = 60
  function_response_types            = ["ReportBatchItemFailures"]
}

module "sqs_lambda" {
  source = "../../modules/lambda"
  bucket = "user-async-lambda"
  function = {
    name = "user-async-lambda"
  }
  env = {
    "USER_BUCKET_NAME" = aws_s3_bucket.user.bucket
    "S3_ENDPOINT"      = "http://s3.localhost.localstack.cloud:4566"
  }
}

resource "aws_s3_bucket" "user" {
  bucket = "users"
}

resource "aws_api_gateway_rest_api" "my_api" {
  name = "my-api"
  endpoint_configuration {
    types = ["REGIONAL"]
  }
  body = jsonencode({
    openapi = "3.0.1"
    info = {
      title   = "My API Gateway"
      version = "1.0"
    }
    paths = {
      "/users" = {
        post = {
          operationId = "createUser"
          requestBody = {
            content = {
              "application/json" = {
                schema = {
                  type = "object"
                  properties = {
                    name       = { type = "string" },
                    birth_date = { type = "string" }
                  }
                }
              }
            }
            required = true
          }
          x-amazon-apigateway-integration = {
            httpMethod           = "POST"
            payloadFormatVersion = "1.0"
            type                 = "AWS_PROXY"
            uri                  = "${module.localstack-demo-lambda.invoke_arn}"
          }
          responses = {
            "200" = {
              description = "200 response"
              content = {
                "application/json" = {
                  schema = {
                    type = "object"
                    properties = {
                      id = { type = "string" }
                    }
                  }
                }
              }
            }
          }
        }
        get = {
          operationId = "listUsers"
          x-amazon-apigateway-integration = {
            httpMethod           = "POST"
            payloadFormatVersion = "1.0"
            type                 = "AWS_PROXY"
            uri                  = "${module.localstack-demo-lambda.invoke_arn}"
          }
          responses = {
            "200" = {
              description = "200 response"
              content = {
                "application/json" = {
                  schema = {
                    type = "array"
                    items = {
                      type = "object"
                      properties = {
                        id   = { type = "string" }
                        name = { type = "string" }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  })
}

resource "aws_api_gateway_deployment" "deployment" {
  triggers = {
    redeployment = sha1(jsonencode(aws_api_gateway_rest_api.my_api.body))
  }
  rest_api_id = aws_api_gateway_rest_api.my_api.id
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_stage" "dev" {
  deployment_id = aws_api_gateway_deployment.deployment.id
  rest_api_id   = aws_api_gateway_rest_api.my_api.id
  stage_name    = "dev"
}

