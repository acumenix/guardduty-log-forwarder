# GuardDuty Log Forwarder to S3

`guardduty-log-forwarder` is a log forwarder to S3 bucket from CloudWatch Event of GuardDuty. This tool is implemented in Go and based on AWS Lambda as Serverless Application by AWS CloudFormation (AWS Serverless Application Model).

## Quick Start

### Prerequsite

- [aws-cil](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-welcome.html)
- AWS resources
  - Credentials: This must be able to deploy CloudFormation (e.g. AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY) from aws-cli.
  - S3 bucket: As GuardDuty log store

### Deploy

```
$ curl -O https://s3-ap-northeast-1.amazonaws.com/guardduty-log-forwarder/templates/latest.yml
$ aws cloudformation deploy \
      --template-file sam.yml \
      --stack-name <YOUR_CLOUD_FORMATION_STACK_NAME> \
      --capabilities CAPABILITY_IAM \
      --parameter-overrides DstS3Bucket=<DST_S3_BUCKET_NAME> DstS3Region=<DST_S3_REGION>
```

You can specify following parameters in `--parameter-overrides` option:

- `LambdaRoleArn`: optional, IAM role ARN for Lambda. A new IAM role will be deployed if not specifing.
- `DstS3Bucket`: **required**, S3 Bucket name to store GuardDuty logs.
- `DstS3Prefix`: optional, S3 Key prefix to store GuardDuty logs. Default prefix is empty. NOTE: `/` is not automatically completed.
- `DstS3Region`: **required**, AWS region name such as `ap-northeast-1` of S3 Bucket to store GuardDuty logs.


## Development

### Setup

```bash
$ git clone git@github.com:m-mizutani/guardduty-log-forwarder.git
$ cd guardduty-log-forwarder
$ dep ensure
```

### Deploy

Create a config file like following as `param.cfg`

```conf
StackName=guardduty-log-forwarder
CodeS3Bucket=my-security-logs.mgmt
CodeS3Prefix=functions

RoleArn=arn:aws:iam::1234567890:role/LambdaGuardDutyLogForwarder
S3Bucket=my-security-logs
S3Prefix=guardduty/
S3Region=ap-northeast-1
```

Then, invoke deploy command.

```bash
$ make
```

### Test

Prepare a test preference file `testpref.json` at top of repository. (a.k.a. in `guardduty-log-forwarder` folder )

```json
{
    "s3_bucket": "mizutani-test",
    "s3_Prefix": "guardduty-forwarder",
    "s3_region": "ap-northeast-1"
}
```

Then, invoke test command.

```
$ go test -v
```
