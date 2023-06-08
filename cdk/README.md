# Prerequisite

Make sure you have Go v1.18.0 or higher.

# Welcome to your CDK Go project!

This is a K2 project for CDK development with Go.

The `cdk.json` file tells the CDK toolkit how to execute your app.

## Useful commands

 * `cdk deploy`      deploy this stack to your default AWS account/region
 * `cdk diff`        compare deployed stack with current state
 * `cdk synth`       emits the synthesized CloudFormation template
 * `go test`         run unit tests

## Coding conventions

Stack naming conventions:
* Stack name for dev environment must follow the pattern `k2-dev*`
* Stack name for prod environment must follow the pattern `k2-prod*`
> Before deploying, make sure K2_ENV_NAME variable is set to any of "k2-dev" or "k2-prod"

To set AWS profile variable on Windows, run powershell command:
```
$Env:AWS_PROFILE='k2-prod'
$Env:AWS_PROFILE='k2-dev'
```

To update only one stack without dependencies:
```
cdk deploy <stack> --exclusively
cdk diff <stack> --exclusively
```
