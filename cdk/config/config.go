package config

import (
	"strings"

	"github.com/aws/jsii-runtime-go"
)

type Environment string

const (
	devEnvironment  Environment = "k2-dev"
	prodEnvironment Environment = "k2-prod"
)

const (
	devAccount       = "041584911022"
	devRegion        = "us-east-1"
	devAllocatedCIDR = "10.130.144.0/21"
)

const (
	prodAccount       = "133360657404"
	prodRegion        = "us-east-1"
	prodAllocatedCIDR = "10.200.64.0/21"
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
