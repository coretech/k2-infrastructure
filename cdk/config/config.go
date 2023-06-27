package config

import (
	"strings"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"

	"github.com/aws/aws-cdk-go/awscdk/v2/awseks"

	"github.com/aws/jsii-runtime-go"
)

type Environment string

const (
	devEnvironment  Environment = "k2-dev"
	prodEnvironment Environment = "k2-prod"
)

const (
	IDTNetworkCIDR = "10.0.0.0/8"
	VpnClientsCIDR = "169.132.0.0/16"
)

// dev related configurations.
const (
	devAccount                               = "041584911022"
	devRegion                                = "us-east-1"
	devAllocatedCIDR                         = "10.130.144.0/21"
	devBrandMediaKitCloudFrontCertificateARN = "arn:aws:acm:us-east-1:041584911022:certificate/4f213c01-440d-4563-85e5-cf141248e849"
	devBrandMediaKitDomainName               = "dev.brand-media-kit.idt.net"
	devTransitGatewayID                      = "tgw-036d427287661ff13"
)

// prod related configurations.
const (
	prodAccount                               = "133360657404"
	prodRegion                                = "us-east-1"
	prodAllocatedCIDR                         = "10.200.64.0/21"
	prodBrandMediaKitCloudFrontCertificateARN = "arn:aws:acm:us-east-1:133360657404:certificate/8b482f77-3907-43b1-980f-c4e4e544bf9d"
	prodBrandMediaKitDomainName               = "prod.brand-media-kit.idt.net"
	prodTransitGatewayID                      = "tgw-046866a4b96cbede2"
)

// Parse parses the environment to deploy infrastructure to from K2_ENV_NAME environment variable.
// If no such variable is set, Parse will return k2-dev.
func Parse(envName string) Environment {
	if strings.EqualFold(envName, string(prodEnvironment)) {
		return prodEnvironment
	}
	return devEnvironment
}

func (e Environment) GetName() *string {
	switch e {
	case devEnvironment:
		return jsii.String(string(devEnvironment))
	case prodEnvironment:
		return jsii.String(string(prodEnvironment))
	default:
		return nil
	}
}

func (e Environment) GetAccount() *string {
	switch e {
	case devEnvironment:
		return jsii.String(devAccount)
	case prodEnvironment:
		return jsii.String(prodAccount)
	default:
		return nil
	}
}

func (e Environment) GetRegion() *string {
	switch e {
	case devEnvironment:
		return jsii.String(devRegion)
	case prodEnvironment:
		return jsii.String(prodRegion)
	default:
		return nil
	}
}

func (e Environment) GetTags() *map[string]*string {
	tags := make(map[string]*string)

	switch e {
	case devEnvironment:
		tags["tech:lob"] = jsii.String("tech_dev_imtu")
		tags["tech:team_name"] = jsii.String("team_mtuoam")
		tags["tech:environment_type"] = jsii.String("dev")
		tags["tech:application_group"] = jsii.String("mtu-services")
	case prodEnvironment:
		tags["tech:lob"] = jsii.String("tech_dev_imtu")
		tags["tech:team_name"] = jsii.String("team_mtuoam")
		tags["tech:environment_type"] = jsii.String("prod")
		tags["tech:application_group"] = jsii.String("mtu-services")
	default:
	}

	return &tags
}

func (e Environment) GetAllocatedCIDR() *string {
	switch e {
	case devEnvironment:
		return jsii.String(devAllocatedCIDR)
	case prodEnvironment:
		return jsii.String(prodAllocatedCIDR)
	default:
		return nil
	}
}

// GetTransitGatewayID fetches Transit Gateway ID (Transit Gateways are provisioned by NETENG).
// Transit Gateway for MTUOAM Prod - https://idtjira.atlassian.net/browse/NETENG-6239.
func (e Environment) GetTransitGatewayID() *string {
	switch e {
	case devEnvironment:
		return jsii.String(devTransitGatewayID)
	case prodEnvironment:
		return jsii.String(prodTransitGatewayID)
	default:
		return nil
	}
}

func (e Environment) GetBrandMediaKitCloudFrontCertificateARN() *string {
	switch e {
	case devEnvironment:
		return jsii.String(devBrandMediaKitCloudFrontCertificateARN)
	case prodEnvironment:
		return jsii.String(prodBrandMediaKitCloudFrontCertificateARN)
	default:
		return nil
	}
}

// DNS records are being handled in IDT Prod account - https://idtjira.atlassian.net/browse/SYS-21972.
// For certificates - go to AWS Certificate Manager (ACM).
func (e Environment) GetBrandMediaKitDomainNames() *[]*string {
	domainNames := make([]*string, 0)
	switch e {
	case devEnvironment:
		domainNames = append(domainNames, jsii.String(devBrandMediaKitDomainName))
		return &domainNames
	case prodEnvironment:
		domainNames = append(domainNames, jsii.String(prodBrandMediaKitDomainName))
		return &domainNames
	default:
		return nil
	}
}

func (e Environment) GetGenericNodeGroupOptionsForEKS() *awseks.NodegroupOptions {
	instanceTypes := make([]awsec2.InstanceType, 0)
	labels := map[string]*string{
		"node_group":    jsii.String("system"),
		"workload":      jsii.String("generic"),
		"ingress":       jsii.String("traefik"),
		"services_type": jsii.String("stateless"),
	}

	switch e {
	case devEnvironment:
		instanceTypes = append(instanceTypes, awsec2.NewInstanceType(jsii.String("m5.large")))
		return &awseks.NodegroupOptions{
			DesiredSize:   jsii.Number[float64](2),
			InstanceTypes: &instanceTypes,
			Labels:        &labels,
			MaxSize:       jsii.Number[float64](2),
		}
	case prodEnvironment:
		instanceTypes = append(instanceTypes, awsec2.NewInstanceType(jsii.String("m5.xlarge")))
		return &awseks.NodegroupOptions{
			DesiredSize:   jsii.Number[float64](2),
			InstanceTypes: &instanceTypes,
			Labels:        &labels,
			MaxSize:       jsii.Number[float64](10),
		}
	default:
		return nil
	}
}
