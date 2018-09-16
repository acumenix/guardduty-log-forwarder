# GuardDuty Log Uploader to S3

`guardduty-log-uploader` is a log uploader to S3 bucket from CloudWatch Event of GuardDuty. This tool is implemented in Go and based on AWS Lambda as Serverless Application.

## Prerequsite

- Go >= 1.10.3
- [mage](https://github.com/magefile/mage) >= 2 (Go based task runner)
- AWS resources
  - Credentials that can deploy CloudFormation (e.g. AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY)
  - IAM Role for Lambda
  - S3 bucket(s)
    - For uploading excutable binary
    - As log repository

## Usage

### Setup

```bash
$ git clone git@github.com:m-mizutani/guardduty-log-uploader.git
$ cd guardduty-log-uploader
$ dep ensure
$ mage build
```

### Deploy

Create a config file like following as `myconf.cfg`

```conf
StackName=guardduty-log-uploader
CodeS3Bucket=my-security-logs.mgmt
CodeS3Prefix=functions
Region=ap-northeast-1

RoleArn=arn:aws:iam::1234567890:role/LambdaGuardDutyLogUploader
S3Bucket=my-security-logs
S3Prefix=guardduty/
S3Region=ap-northeast-1
```

Then, invoke deploy command.

```bash
$ env PARAM_FILE=myconf.cfg mage deploy
```

### Test

Prepare a test preference file `testpref.json` at top of repository. (a.k.a. in `guardduty-log-uploader` folder )

```json
{
    "s3_bucket": "mizutani-test",
    "s3_Prefix": "guardduty-uploader",
    "s3_region": "ap-northeast-1"
}
```

Then, invoke test command.

```
$ mage test
```