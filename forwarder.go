package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

// Event is data from CloudWatch event of GuardDuty
type Event struct {
	Version    string      `json:"version"`
	ID         string      `json:"id"`
	DetailType string      `json:"detail-type"`
	Source     string      `json:"source"`
	Account    string      `json:"account"`
	Time       string      `json:"time"`
	Region     string      `json:"region"`
	Resources  interface{} `json:"resources"`
	Detail     interface{} `json:"detail"`
}

func main() {
	lambda.Start(handleRequest)
}

func getRegion(ctx context.Context) string {
	lambdaCtx, ok := lambdacontext.FromContext(ctx)
	if !ok {
		panic("can not convert to lambdaCtx")
	}

	arn := strings.Split(lambdaCtx.InvokedFunctionArn, ":")
	return arn[3]
}

func handleRequest(ctx context.Context, event Event) (string, error) {
	args := Args{
		Event:    event,
		S3Bucket: os.Getenv("S3_BUCKET"),
		S3Prefix: os.Getenv("S3_PREFIX"),
		S3Region: os.Getenv("S3_REGION"),
	}

	return Handler(&args)
}

type Args struct {
	Event    Event
	S3Bucket string
	S3Prefix string
	S3Region string
}

func Handler(args *Args) (string, error) {
	log.Println(args)

	jdata, err := json.Marshal(&args.Event)
	if err != nil {
		return "", errors.Wrap(err, "Fail to convert CloudWatch Event")
	}

	s3Prefix := args.S3Prefix
	if !strings.HasSuffix(s3Prefix, "/") {
		s3Prefix += "/"
	}

	dt, err := time.Parse("2006-01-02T15:04:05Z", args.Event.Time)
	s3Key := fmt.Sprintf("%s%s/%04d/%02d/%02d/%02d/%s.json",
		s3Prefix, args.Event.Region, dt.Year(), dt.Month(), dt.Day(),
		dt.Hour(), args.Event.ID)

	ssn := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(args.S3Region),
	}))

	client := s3.New(ssn)
	resp, err := client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(args.S3Bucket),
		Key:    aws.String(s3Key),
		Body:   aws.ReadSeekCloser(bytes.NewReader(jdata)),
	})
	if err != nil {
		log.Println(err.Error())
	}

	log.Println(awsutil.StringValue(resp))

	return s3Key, nil
}
