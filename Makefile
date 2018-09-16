CONFIG_FILE ?= "param.cfg"

CODE_S3_BUCKET := $(shell cat $(CONFIG_FILE) | grep CodeS3Bucket | cut -d = -f 2)
CODE_S3_PREFIX := $(shell cat $(CONFIG_FILE) | grep CodeS3Prefix | cut -d = -f 2)
STACK_NAME := $(shell cat $(CONFIG_FILE) | grep StackName | cut -d = -f 2)
PARAMETERS := $(shell cat $(CONFIG_FILE) | grep -e LambdaRoleArn -e DstS3Bucket -e DstS3Prefix -e DstS3Region | tr '\n' ' ')
TEMPLATE_FILE=template.yml

all: deploy

build/forwarder: *.go
	env GOARCH=amd64 GOOS=linux go build -o build/forwarder

test:
	go test -v

sam.yml: build/forwarder template.yml
	aws cloudformation package \
		--template-file $(TEMPLATE_FILE) \
		--s3-bucket $(CODE_S3_BUCKET) \
		--s3-prefix $(CODE_S3_PREFIX) \
		--output-template-file sam.yml

deploy: sam.yml
	aws cloudformation deploy \
		--template-file sam.yml \
		--stack-name $(STACK_NAME) \
		--capabilities CAPABILITY_IAM \
		--parameter-overrides $(PARAMETERS)
