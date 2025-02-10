## localstack-demo
This is a demo project to show how to use localstack with Golang and Terraform.

## Architecture
<img width="821" alt="Screenshot 2025-02-10 at 18 12 37" src="https://github.com/user-attachments/assets/cc692de9-9b41-48fe-8535-230231220539" />

## Set up
### 1. clone this repository
### 2. install localstack
create an account at https://www.localstack.cloud/ and install.
```bash
brew install localstack/tap/localstack-cli
```
### 3. download localstack desktop from [here](https://hub.docker.com/extensions/localstack/localstack-docker-desktop)
you can see emulated aws resources on it or at https://app.localstack.cloud/

### 4. setup terraform
```
brew tap hashicorp/tap
brew install hashicorp/tap/terraform
```

### 5. setup your local AWS profile
add `localstack` profile
```bash
yuji.tamura@MacBook-Pro projects % cat ~/.aws/config 
[profile localstack]
region = eu-west-2
output = json

yuji.tamura@MacBook-Pro projects % cat ~/.aws/credentials 
[localstack]
aws_access_key_id = test
aws_secret_access_key = test
```
### 6. run LocalStack on docker
```bash
docker-compose up
```

### 7. deploy using terraform
```
yuji.tamura@MacBook-Pro localstack-demo % cd ./terraform/environments/local 
yuji.tamura@MacBook-Pro local % AWS_PROFILE=localstack terraform init
Initializing the backend...
Initializing modules...
Initializing provider plugins...
- Reusing previous version of hashicorp/aws from the dependency lock file
- Reusing previous version of hashicorp/null from the dependency lock file
- Using previously-installed hashicorp/null v3.2.3
- Using previously-installed hashicorp/aws v5.86.0

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
yuji.tamura@MacBook-Pro local % AWS_PROFILE=localstack terraform plan
yuji.tamura@MacBook-Pro local % AWS_PROFILE=localstack terraform apply
```

## Useful commands
- send a request to API gateway

  ```bash
  curl -X POST --header "Content-Type: application/json" \
  --data '{"name":"xyz","birth_date":"1995-02-21"}' \
  http://<<api gateway's api id>>.execute-api.localhost.localstack.cloud:4566/dev/users
  ```
  note: `<<api gateway's api id>>` can be found on LocalStack Desktop or on https://app.localstack.cloud/inst/default/resources/apigateway/
  
- send a message to sqs

  ```bash
  AWS_PROFILE=localstack aws --endpoint-url=http://localhost:4566 sqs send-message \
   --queue-url=http://sqs.eu-west-2.localhost.localstack.cloud:4566/000000000000/test \
   --message-body='{"id": "abc","name":"yuji","birth_date":"1995-02-21"}'
  ```

