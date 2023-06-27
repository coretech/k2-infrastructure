package stacks

import (
	"cdk/config"
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"

	"github.com/aws/aws-cdk-go/awscdk/v2/lambdalayerkubectl"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseks"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const (
	eksStackID = "%s-eks"
)

type EKSStackProps struct {
	awscdk.StackProps
	config.Environment
	awsec2.Vpc
}

func NewEKSStack(scope constructs.Construct, props *EKSStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, jsii.String(fmt.Sprintf(eksStackID, *props.Environment.GetName())), &sprops)
	vpcSubnets := []*awsec2.SubnetSelection{{SubnetGroupName: jsii.String(defaultApplicationSubnetGroupName)}}

	// Define the EKS cluster
	cluster := awseks.NewCluster(stack, jsii.String("eks"), &awseks.ClusterProps{
		Version:     awseks.KubernetesVersion_V1_26(),
		ClusterName: jsii.String(fmt.Sprintf(eksStackID, *props.Environment.GetName())),
		Vpc:         props.Vpc,
		VpcSubnets:  &vpcSubnets,
		AlbController: &awseks.AlbControllerOptions{
			Version: awseks.AlbControllerVersion_V2_5_1(),
		},
		KubectlLayer:    lambdalayerkubectl.NewKubectlLayer(stack, jsii.String("kubectl-layer")),
		DefaultCapacity: jsii.Number[float64](0),
	})

	//awscdk.NewCfnOutput(scope, jsii.String("ClusterSecurityGroup"), &awscdk.CfnOutputProps{
	//	Value:       cluster.ClusterSecurityGroupId(),
	//	Description: jsii.String("security group for the cluster"),
	//	ExportName:  jsii.String("security-group"),
	//})

	// Route 53 permissions (external DNS)
	policyStatements := []awsiam.PolicyStatement{
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			Actions:   jsii.Strings("route53:GetChange"),
			Resources: jsii.Strings("arn:aws:route53:::change/*"),
		}),
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			Actions:   jsii.Strings("route53:ChangeResourceRecordSets"),
			Resources: jsii.Strings("arn:aws:route53:::hostedzone/*"),
		}),
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			Actions:   jsii.Strings("route53:ListHostedZones", "route53:ListResourceRecordSets"),
			Resources: jsii.Strings("*"),
		}),
	}

	route53Policy := awsiam.NewPolicy(stack, jsii.String("route-53-policy"), &awsiam.PolicyProps{
		Statements: &policyStatements,
	})

	externalDnsSa := cluster.AddServiceAccount(jsii.String("external-dns"), &awseks.ServiceAccountOptions{
		Name:      jsii.String("external-dns"),
		Namespace: jsii.String("kube-system"),
	})

	externalDnsSa.Role().AttachInlinePolicy(route53Policy)

	manifest := map[string]any{
		"apiVersion": "v1",
		"kind":       "Namespace",
		"metadata": map[string]string{
			"name": "cert-manager",
		},
	}

	certManagerNamespace := cluster.AddManifest(jsii.String("cert-manager-namespace"), &manifest)

	certManagerSa := cluster.AddServiceAccount(jsii.String("cert-manager"), &awseks.ServiceAccountOptions{
		Name:      jsii.String("cert-manager"),
		Namespace: jsii.String("cert-manager"),
	})

	certManagerSa.Role().AttachInlinePolicy(route53Policy)
	certManagerSa.Node().AddDependency(certManagerNamespace)

	cluster.ClusterSecurityGroup().AddIngressRule(
		awsec2.Peer_Ipv4(jsii.String(config.VpnClientsCIDR)),
		awsec2.Port_AllIcmp(),
		jsii.String("allow ICMP traffic for VPN clients"),
		jsii.Bool(false),
	)

	// Several node types e.g. generic/data/compute with scaling policies.
	cluster.AddNodegroupCapacity(jsii.String("generic-v1"), props.Environment.GetGenericNodeGroupOptionsForEKS())

	return stack
}
