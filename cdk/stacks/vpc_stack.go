package stacks

import (
	"cdk/config"
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const (
	vpcStackID                        = "%s-vpc"
	defaultApplicationSubnetGroupName = "private"
	publicSubnetGroupName             = "public"
)

type VPCStackProps struct {
	awscdk.StackProps
	config.Environment
}

// MTUOAM Dev allocated CIDR: 10.130.144.0/21 Ref: https://idtjira.atlassian.net/browse/NAT-4322
// MTUOAM Prod allocated CIDR: 10.200.64.0/21 Ref: https://idtjira.atlassian.net/browse/NAT-4323
// https://docs.netgate.com/pfsense/en/latest/network/cidr.html
func NewVPCStack(scope constructs.Construct, props *VPCStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, jsii.String(fmt.Sprintf(vpcStackID, *props.Environment.GetName())), &sprops)
	subnetConfigurations := make([]*awsec2.SubnetConfiguration, 0)

	subnetConfigurations = append(
		subnetConfigurations,
		&awsec2.SubnetConfiguration{
			CidrMask:   jsii.Number[int](24),
			Name:       jsii.String(defaultApplicationSubnetGroupName),
			SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
		},
		&awsec2.SubnetConfiguration{
			CidrMask:   jsii.Number[int](26),
			Name:       jsii.String(publicSubnetGroupName),
			SubnetType: awsec2.SubnetType_PUBLIC,
		},
	)

	// Define the VPC
	awsec2.NewVpc(stack, jsii.String("default"), &awsec2.VpcProps{
		IpAddresses:         awsec2.IpAddresses_Cidr(props.Environment.GetAllocatedCIDR()),
		MaxAzs:              jsii.Number[int](3),
		NatGateways:         jsii.Number[int](3),
		SubnetConfiguration: &subnetConfigurations,
	})

	return stack
}
