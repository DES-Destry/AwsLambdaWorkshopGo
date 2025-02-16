# Serverless Patterns for Go
[Workshop Studio](https://catalog.workshops.aws/serverless-patterns/en-US)

### Access to AWS

First of all download AWS CLI. Just google it.

To execute deploy commands below you have to grant access to AWS to your AWS CLI

Create a pair of Access Key and Secret Key in your AWS Dashboard in AWS IAM. Create a new user, and then you'll be able to create keys for him.

Then run command (this will require Access & Secret Keys):

```bash
aws configure --profile workshop.serverless_patterns
```

### Deploy

First create an AWS Lambda `AWSLambdaDefaultRole` Role with the following policy named `AWSLambdaBasicExecutionRole`.

There's no Go runtime in AWS Lambda, so we have to build our Go code and deploy it as a binary with Amazon Linux 2023 Lambda environment.

Build your Go package:
```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap -tags lambda.norpc main.go
```

Zip the binary:
```bash
zip myFunction.zip bootstrap
```

Publish the function:
```bash
AWS_PROFILE=workshop.serverless_patterns aws lambda create-function --function-name TestGo \
--region us-east-1 \
--runtime provided.al2023 --handler bootstrap \
--role arn:aws:iam::<YOUR_AWS_ACCOUNT>:role/service-role/AWSLambdaDefaultRole \
--zip-file fileb://myFunction.zip
```

More info at: [AWS Lambda: Building Lambda functions with Go](https://docs.aws.amazon.com/lambda/latest/dg/lambda-golang.html)