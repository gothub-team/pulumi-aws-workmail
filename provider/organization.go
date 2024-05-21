package provider

import (
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/workmail"
	"github.com/aws/aws-sdk-go-v2/service/workmail/types"
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
type Organization struct{}

// Each resource has an input struct, defining what arguments it accepts.
type OrganizationArgs struct {
	// The organization alias.
	Alias string `pulumi:"alias"`
	// The email domains to associate with the organization.
	Domains []DomainArgs `pulumi:"domains"`
	// The idempotency token associated with the request.
	ClientToken *string `pulumi:"clientToken"`
	// The AWS Directory Service directory ID.
	DirectoryId *string `pulumi:"directoryId"`
	// The Amazon Resource Name (ARN) of a customer managed key from AWS KMS.
	KmsKeyArn *string `pulumi:"kmsKeyArn"`
	// When true , allows organization interoperability between WorkMail and Microsoft Exchange. If true , you must include a AD Connector directory ID in the request.
	EnableInteroperability bool `pulumi:"enableInteroperability"`
}

type DomainArgs struct {
	// The domain name.
	DomainName string `pulumi:"domainName"`
	// The hosted zone id for the domain.
	HostedZoneId string `pulumi:"hostedZoneId"`
}

// Each resource has a state, describing the fields that exist on the created resource.
type OrganizationState struct {
	// It is generally a good idea to embed args in outputs, but it isn't strictly necessary.
	OrganizationArgs
}

// All resources must implement Create at a minimum.
func (Organization) Create(ctx p.Context, name string, input OrganizationArgs, preview bool) (string, OrganizationState, error) {
	state := OrganizationState{OrganizationArgs: input}
	if preview {
		return name, state, nil
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", state, err
	}

	// Create the WorkMail service client using the config
	workmailclient := workmail.NewFromConfig(cfg)

	organization, err := workmailclient.CreateOrganization(ctx, &workmail.CreateOrganizationInput{
		Alias: &input.Alias,
		Domains: Map(func(domain DomainArgs) types.Domain {
			return types.Domain{
				DomainName:   &domain.DomainName,
				HostedZoneId: &domain.HostedZoneId,
			}
		})(input.Domains),
		ClientToken:            input.ClientToken,
		DirectoryId:            input.DirectoryId,
		KmsKeyArn:              input.KmsKeyArn,
		EnableInteroperability: true,
	})
	if err != nil {
		return "", state, err
	}

	return *organization.OrganizationId, state, nil
}

func Map[T any, U any](f func(T) U) func([]T) []U {
	return func(slice []T) []U {
		newSlice := make([]U, len(slice))
		for i, value := range slice {
			newSlice[i] = f(value)
		}
		return newSlice
	}
}
