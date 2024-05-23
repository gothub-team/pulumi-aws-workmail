package provider

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
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
	// The AWS Region. TODO: This should be passed as a pulumi.Provider
	Region string `pulumi:"region"`
	// The organization alias.
	Alias string `pulumi:"alias"`
	// The domain name.
	DomainName string `pulumi:"domainName"`
	// The hosted zone id for the domain.
	HostedZoneId string `pulumi:"hostedZoneId"`
	// The idempotency token associated with the request.
	ClientToken *string `pulumi:"clientToken,optional"`
	// The AWS Directory Service directory ID.
	DirectoryId *string `pulumi:"directoryId,optional"`
	// The Amazon Resource Name (ARN) of a customer managed key from AWS KMS.
	KmsKeyArn *string `pulumi:"kmsKeyArn,optional"`
	// When true , allows organization interoperability between WorkMail and Microsoft Exchange. If true , you must include a AD Connector directory ID in the request.
	EnableInteroperability *bool `pulumi:"enableInteroperability,optional"`
}

// Each resource has a state, describing the fields that exist on the created resource.
type OrganizationState struct {
	// It is generally a good idea to embed args in outputs, but it isn't strictly necessary.
	OrganizationArgs

	// The organization id.
	OrganizationId string `pulumi:"organizationId"`

	// Mail domain records.
	Records []DnsRecord `pulumi:"records"`
}

type DnsRecord struct {
	Type     string `pulumi:"type"`
	Hostname string `pulumi:"hostname"`
	Value    string `pulumi:"value"`
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
	cfg.Region = input.Region

	// Create the WorkMail service client using the config
	workmailclient := workmail.NewFromConfig(cfg)
	sesclient := ses.NewFromConfig(cfg)

	// Create the organization
	organization, err := workmailclient.CreateOrganization(ctx, &workmail.CreateOrganizationInput{
		Alias: &input.Alias,
		Domains: []types.Domain{
			{
				DomainName:   &input.DomainName,
				HostedZoneId: &input.HostedZoneId,
			},
		},
		ClientToken:            input.ClientToken,
		DirectoryId:            input.DirectoryId,
		KmsKeyArn:              input.KmsKeyArn,
		EnableInteroperability: ifNotNil(input.EnableInteroperability, false),
	})
	if err != nil {
		return "", state, err
	}
	state.OrganizationId = *organization.OrganizationId

	// Wait for the organization to be created
	for {
		org, err := workmailclient.DescribeOrganization(ctx, &workmail.DescribeOrganizationInput{
			OrganizationId: &state.OrganizationId,
		})
		if err != nil {
			return "", state, err
		}

		identityVerification, err := sesclient.GetIdentityVerificationAttributes(ctx, &ses.GetIdentityVerificationAttributesInput{
			Identities: []string{input.DomainName},
		})
		if err != nil {
			return "", state, err
		}
		fmt.Println(*org.State, identityVerification.VerificationAttributes[input.DomainName].VerificationStatus)
		if *org.State == "Active" && identityVerification.VerificationAttributes[input.DomainName].VerificationStatus == "Success" {
			break
		}
		time.Sleep(5 * time.Second)
	}

	_, err = workmailclient.UpdateDefaultMailDomain(ctx, &workmail.UpdateDefaultMailDomainInput{
		OrganizationId: &state.OrganizationId,
		DomainName:     &input.DomainName,
	})
	if err != nil {
		return "", state, err
	}

	mailDomain, err := workmailclient.GetMailDomain(ctx, &workmail.GetMailDomainInput{
		OrganizationId: &state.OrganizationId,
		DomainName:     &input.DomainName,
	})
	if err != nil {
		return "", state, err
	}

	state.Records = Map(func(record types.DnsRecord) DnsRecord {
		return DnsRecord{
			Type:     *record.Type,
			Hostname: *record.Hostname,
			Value:    *record.Value,
		}
	})(mailDomain.Records)

	return state.OrganizationId, state, nil
}

func ifNotNil[T any](ptr *T, def T) T {
	if ptr != nil {
		return *ptr
	}
	return def
}

// The Delete method will run when the resource is deleted.
func (Organization) Delete(ctx p.Context, id string, props OrganizationState) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	cfg.Region = props.Region

	// Create the WorkMail service client using the config
	workmailclient := workmail.NewFromConfig(cfg)
	sesclient := ses.NewFromConfig(cfg)

	organization, err := workmailclient.DescribeOrganization(ctx, &workmail.DescribeOrganizationInput{
		OrganizationId: &id,
	})
	if err != nil {
		return nil
	}

	if *organization.State != "Deleted" {
		// Delete the organization
		_, err = workmailclient.DeleteOrganization(ctx, &workmail.DeleteOrganizationInput{
			OrganizationId:  &id,
			DeleteDirectory: true,
			ForceDelete:     true,
		})
		if err != nil {
			return err
		}

		for {
			organization, err := workmailclient.DescribeOrganization(ctx, &workmail.DescribeOrganizationInput{
				OrganizationId: &id,
			})
			if err != nil {
				return err
			}
			if *organization.State == "Deleted" {
				break
			}
			time.Sleep(5 * time.Second)
		}
	}

	_, err = sesclient.DeleteIdentity(ctx, &ses.DeleteIdentityInput{
		Identity: &props.DomainName,
	})

	return err
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
