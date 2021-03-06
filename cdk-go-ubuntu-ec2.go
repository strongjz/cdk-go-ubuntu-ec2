package main

import (
	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsec2"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type CdkGoUbuntuEc2StackProps struct {
	awscdk.StackProps
}

func NewCdkGoUbuntuEc2Stack(scope constructs.Construct, id string, props *CdkGoUbuntuEc2StackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	subnetConfig := make([]*awsec2.SubnetConfiguration,6)

	s1 := awsec2.SubnetConfiguration{
		Name:       jsii.String("Public-1"),
		SubnetType: awsec2.SubnetType_PUBLIC,
		CidrMask:   jsii.Number(24),
	}
	s2 := awsec2.SubnetConfiguration{
		Name:       jsii.String("Public-2"),
		SubnetType: awsec2.SubnetType_PUBLIC,
		CidrMask:   jsii.Number(24),
	}
	s3 := awsec2.SubnetConfiguration{
		Name:       jsii.String("Public-3"),
		SubnetType: awsec2.SubnetType_PUBLIC,
		CidrMask:   jsii.Number(24),
	}
	s4 := awsec2.SubnetConfiguration{
		Name:       jsii.String("Private-1"),
		SubnetType: awsec2.SubnetType_PRIVATE,
		CidrMask:   jsii.Number(24),
	}
	s5 := awsec2.SubnetConfiguration{
		Name:       jsii.String("Private-2"),
		SubnetType: awsec2.SubnetType_PRIVATE,
		CidrMask:   jsii.Number(24),
	}
	s6 := awsec2.SubnetConfiguration{
		Name:       jsii.String("Private-3"),
		SubnetType: awsec2.SubnetType_PRIVATE,
		CidrMask:   jsii.Number(24),
	}

	subnetConfig = append(subnetConfig,&s1,&s2,&s3,&s4,&s5,&s6)

	var vpc = awsec2.NewVpc(stack,jsii.String("CdkGoUbuntuEc2Stack"), &awsec2.VpcProps{
		Cidr:                   jsii.String("10.0.0.0/16"),
		DefaultInstanceTenancy: awsec2.DefaultInstanceTenancy_DEFAULT,
		EnableDnsHostnames:     jsii.Bool(true),
		EnableDnsSupport:       jsii.Bool(true),
		MaxAzs:                 jsii.Number(3),
		NatGatewayProvider:     awsec2.NatProvider_Gateway(&awsec2.NatGatewayProps{}),
		NatGateways:            jsii.Number(1),
		SubnetConfiguration:    &subnetConfig,
	})


	imageMap := make(map[string]*string)

	var ubuntu = "ami-019212a8baeffb0fa"

	imageMap["us-east-1"] = jsii.String(ubuntu)
	
	var script = "apt-get update -y \n " +
		"apt-get install -y git awscli ec2-instance-connect \n" +
		"until git clone https://github.com/aws-quickstart/quickstart-linux-utilities.git; do echo \"Retrying\"; done\n" +
		"cd /quickstart-linux-utilities\n" +
		"source quickstart-cfn-tools.source\n  " +
		"qs_update-os || qs_err\n" +
		"qs_bootstrap_pip || qs_err\n" +
		"qs_aws-cfn-bootstrap || qs_err\n  " +
		"mkdir -p /opt/aws/bin\n" +
		"ln -s /usr/local/bin/cfn-* /opt/aws/bin/"

	prop := awsec2.GenericLinuxImageProps{
		UserData: awsec2.UserData_Custom(jsii.String(script)),
	}
	ec2ID := "CdkGoUbuntuEc2Stack"
	image := awsec2.NewGenericLinuxImage(&imageMap,&prop)

	instanceProp := awsec2.InstanceProps{
		InstanceType: awsec2.NewInstanceType(jsii.String("X86_64")),
		MachineImage:     image,
		Vpc:              vpc,
		InstanceName:              jsii.String(ec2ID),
		UserData:                  prop.UserData,
		VpcSubnets: &awsec2.SubnetSelection{
			SubnetType:        awsec2.SubnetType_PUBLIC,
		},
	}

	awsec2.NewInstance(stack,&ec2ID,&instanceProp)


	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewCdkGoUbuntuEc2Stack(app, "CdkGoUbuntuEc2Stack", &CdkGoUbuntuEc2StackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	//return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	return &awscdk.Environment{
	  Region:  jsii.String("us-west-2"),
	 }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
