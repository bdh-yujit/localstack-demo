variable "bucket" {
  description = "Lambda S3 Bucket name"
  type        = string
}

variable "function" {
  description = "Lambda function name"
  type = object({
    name = string
  })
}

variable "env" {
  description = "Environment variables"
  type = map(string)
}
