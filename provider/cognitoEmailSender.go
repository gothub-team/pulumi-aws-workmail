package provider

import (
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	p "github.com/pulumi/pulumi-go-provider"
)

// Each resource has a controlling struct.
// Resource behavior is determined by implementing methods on the controlling struct.
// The `Create` method is mandatory, but other methods are optional.
// - Check: Remap inputs before they are typed.
// - Diff: Change how instances of a resource are compared.
// - Update: Mutate a resource in place.
// - Read: Get the state of a resource from the backing provider.
// - Delete: Custom logic when the resource is deleted.
// - Annotate: Describe fields and set defaults for a resource.
// - WireDependencies: Control how outputs and secrets flows through values.
type CognitoEmailSender struct{}

// Each resource has an input struct, defining what arguments it accepts.
type CognitoEmailSenderArgs struct {
	// ID of the cognito user pool that should be updated with the custom email sender.
	UserPoolId string `pulumi:"userPoolId"`
	// Arn of the lambda that is responsible for sending emails
	LambdaArn string `pulumi:"lambdaArn"`
	// Arn of the KMS key that is used to encrypt verification codes sent via email.
	KmsKeyArn string `pulumi:"kmsKeyArn"`
}

// Each resource has a state, describing the fields that exist on the created resource.
type CognitoEmailSenderState struct {
	// It is generally a good idea to embed args in outputs, but it isn't strictly necessary.
	CognitoEmailSenderArgs
}

// All resources must implement Create at a minimum.
func (CognitoEmailSender) Create(ctx p.Context, name string, input CognitoEmailSenderArgs, preview bool) (string, CognitoEmailSenderState, error) {
	state := CognitoEmailSenderState{CognitoEmailSenderArgs: input}
	if preview {
		return name, state, nil
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", state, err
	}

	cognitoclient := cognitoidentityprovider.NewFromConfig(cfg)

	_, err = cognitoclient.UpdateUserPool(ctx, &cognitoidentityprovider.UpdateUserPoolInput{
		UserPoolId: &input.UserPoolId,
		LambdaConfig: &types.LambdaConfigType{
			CustomEmailSender: &types.CustomEmailLambdaVersionConfigType{
				LambdaVersion: "V1_0",
				LambdaArn:     &input.LambdaArn,
			},
			KMSKeyID: &input.KmsKeyArn,
		},
		AutoVerifiedAttributes: []types.VerifiedAttributeType{
			"email",
		},
	})
	if err != nil {
		return "", state, err
	}

	return input.LambdaArn, state, nil
}

func (CognitoEmailSender) Diff(ctx p.Context, id string, olds CognitoEmailSenderState, news CognitoEmailSenderArgs) (p.DiffResponse, error) {
	diffs := make(map[string]p.PropertyDiff)
	hasChanges := false

	if true {
		diffs["lambdaArn"] = p.PropertyDiff{Kind: p.UpdateReplace, InputDiff: true}
		hasChanges = true
	}

	if olds.UserPoolId != news.UserPoolId {
		diffs["userPoolId"] = p.PropertyDiff{Kind: p.UpdateReplace, InputDiff: true}
		hasChanges = true
	}

	if olds.KmsKeyArn != news.KmsKeyArn {
		diffs["kmsKeyArn"] = p.PropertyDiff{Kind: p.UpdateReplace, InputDiff: true}
		hasChanges = true
	}

	return p.DiffResponse{HasChanges: hasChanges, DetailedDiff: diffs, DeleteBeforeReplace: true}, nil
}

// The Delete method will run when the resource is deleted.
func (CognitoEmailSender) Delete(ctx p.Context, id string, props CognitoEmailSenderState) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	cognitoclient := cognitoidentityprovider.NewFromConfig(cfg)

	_, err = cognitoclient.UpdateUserPool(ctx, &cognitoidentityprovider.UpdateUserPoolInput{
		UserPoolId:             &props.UserPoolId,
		LambdaConfig:           nil,
		AutoVerifiedAttributes: []types.VerifiedAttributeType{},
	})

	return err
}
