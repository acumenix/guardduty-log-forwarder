#!/bin/bash

REGIONS=(`cat scripts/guardduty-regions.txt | tr '\n' ' '`)

if [ $# -le 3 ]; then
    echo "usage) $0 StackName DstS3Region DstS3Bucket DstS3Prefix [LambdaArn]"
    exit 1
fi

if [ "$5" != "" ]; then
   LAMBDA_ARN="LambdaArn=$5"
else
   LAMBDA_ARN=""
fi

echo $LAMBDA_ARN

for region in ${REGIONS[@]}; do
    TEMPLATE="sam.$region.yml"
    curl -sf -o $TEMPLATE https://s3-$region.amazonaws.com/cfn-assets.$region/guardduty-log-forwarder/templates/latest.yml
    aws --region $region cloudformation deploy \
        --template-file $TEMPLATE \
        --stack-name $1 \
        --capabilities CAPABILITY_IAM \
        --parameter-overrides DstS3Region=$2 DstS3Bucket=$3 DstS3Prefix=$4 $LAMBDA_ARN
done
