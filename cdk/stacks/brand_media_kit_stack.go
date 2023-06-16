package stacks

import (
	"cdk/config"
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfrontorigins"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const (
	brandMediaKitStackID                = "%s-brand-media-kit"
	cloudFrontDistributionName          = "brand-media-kit-cloudfront"
	cloudFrontOriginAccessIdentityID    = "brand-media-kit-cf-origin-access-identity"
	s3BucketName                        = "brand-media-kit"
	s3BrandMediaKitRelatedDocumentsPath = "/assets"
)

type BrandMediaKitStackProps struct {
	awscdk.StackProps
	config.Environment
}

// Every constuct related to Brand Media Kit should go here
func NewBrandMediaKitStack(scope constructs.Construct, props *BrandMediaKitStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, jsii.String(fmt.Sprintf(brandMediaKitStackID, *props.Environment.GetName())), &sprops)

	// S3 bucket for WordPress instances to store media and assets.
	originBucket := awss3.NewBucket(stack, jsii.String(fmt.Sprintf(brandMediaKitStackID, *props.Environment.GetName())), &awss3.BucketProps{
		BucketName:    jsii.String(fmt.Sprintf(brandMediaKitStackID, *props.Environment.GetName())),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		Versioned:     jsii.Bool(true),
	})

	originAccessIdentity := awscloudfront.NewOriginAccessIdentity(stack, jsii.String(cloudFrontOriginAccessIdentityID), &awscloudfront.OriginAccessIdentityProps{
		Comment: jsii.String("Origin Access Identity used for accessing S3 by CloudFront distribution"),
	})

	awscloudfront.NewDistribution(stack, jsii.String(cloudFrontDistributionName), &awscloudfront.DistributionProps{
		DefaultBehavior: &awscloudfront.BehaviorOptions{
			Origin: awscloudfrontorigins.NewS3Origin(originBucket, &awscloudfrontorigins.S3OriginProps{
				OriginId:             jsii.String(fmt.Sprintf(brandMediaKitStackID, *props.Environment.GetName())),
				OriginPath:           jsii.String(s3BrandMediaKitRelatedDocumentsPath),
				OriginAccessIdentity: originAccessIdentity,
			}),
		},
		Certificate:       nil,
		Comment:           jsii.String("CloudFront Distribution for Brand Media Kit"),
		DefaultRootObject: nil,
		DomainNames:       nil,
		Enabled:           jsii.Bool(true),
		EnableIpv6:        nil,
		EnableLogging:     jsii.Bool(true),
		ErrorResponses:    nil,
	})

	// Grant CloudFront access to the S3 bucket
	originBucket.GrantRead(originAccessIdentity, jsii.String("GrantCloudFrontAccess"))

	return stack
}
