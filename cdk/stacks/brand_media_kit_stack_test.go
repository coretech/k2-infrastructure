package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/jsii-runtime-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBrandMediaKitStackCreation(t *testing.T) {
	app := awscdk.NewApp(nil)
	// Create a testing stack with a mocked Construct scope
	stack := NewBrandMediaKitStack(app, "test", &BrandMediaKitStackProps{
		StackProps: awscdk.StackProps{
			Env: &awscdk.Environment{
				Region: jsii.String("us-east-1"),
			},
		},
	})

	// Get the S3 bucket from the stack and assert its properties
	bucket := awss3.Bucket_FromBucketArn(stack, jsii.String("brand-media-kit"), stack.FormatArn(&awscdk.ArnComponents{
		Service:      jsii.String("s3"),
		Resource:     jsii.String(s3BucketName),
		ResourceName: jsii.String(s3BucketName),
	}))

	assert.NotNil(t, bucket)
	assert.Equal(t, s3BucketName, *bucket.BucketName())
}
