package stacks

import (
	"cdk/config"
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const (
	brandMediaKitStackID = "%s-brand-media-kit"
	s3BucketName         = "brand-media-kit"
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
	awss3.NewBucket(stack, jsii.String(fmt.Sprintf(brandMediaKitStackID, *props.Environment.GetName())), &awss3.BucketProps{
		BucketName:    jsii.String(fmt.Sprintf(brandMediaKitStackID, *props.Environment.GetName())),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		Versioned:     jsii.Bool(true),
	})

	return stack
}
