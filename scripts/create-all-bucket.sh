#!/bin/bash

REGIONS=(`cat scripts/guardduty-regions.txt | tr '\n' ' '`)

for region in ${REGIONS[@]}; do
    echo configure $region
    if [ "$region" = "us-east-1" ]; then
	aws s3api create-bucket --bucket cfn-assets.$region --region $region 
    else
	aws s3api create-bucket --bucket cfn-assets.$region --region $region --create-bucket-configuration LocationConstraint=$region
    fi

    POLICY=`cat scripts/bucket-policy.json | sed -e "s/__REGION__/$region/g"`
    aws s3api put-bucket-policy --bucket cfn-assets.$region --policy "$POLICY"    
done
