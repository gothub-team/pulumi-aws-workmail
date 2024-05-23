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
type DefaultDomain struct{}

// Each resource has an input struct, defining what arguments it accepts.
type DefaultDomainArgs struct {
	// The AWS Region. TODO: This should be passed as a pulumi.Provider
	Region string `pulumi:"region"`
	// The domain name.
	DomainName string `pulumi:"domainName"`
	// The organization the domain should be associated with.
	OrganizationId string `pulumi:"organizationId"`
	// The idempotency token associated with the request.
	ClientToken *string `pulumi:"clientToken,optional"`
}

// Each resource has a state, describing the fields that exist on the created resource.
type DefaultDomainState struct {
	// It is generally a good idea to embed args in outputs, but it isn't strictly necessary.
	DefaultDomainArgs

	// Mail domain records.
	Records []DnsRecord `pulumi:"records"`
}

type DnsRecord struct {
	Type     string `pulumi:"type"`
	Hostname string `pulumi:"hostname"`
	Value    string `pulumi:"value"`
}

// All resources must implement Create at a minimum.
func (DefaultDomain) Create(ctx p.Context, name string, input DefaultDomainArgs, preview bool) (string, DefaultDomainState, error) {
	state := DefaultDomainState{DefaultDomainArgs: input}
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

	_, err = workmailclient.RegisterMailDomain(ctx, &workmail.RegisterMailDomainInput{
		DomainName:     &input.DomainName,
		OrganizationId: &input.OrganizationId,
		ClientToken:    input.ClientToken,
	})
	if err != nil {
		return "", state, err
	}

	_, err = workmailclient.UpdateDefaultMailDomain(ctx, &workmail.UpdateDefaultMailDomainInput{
		OrganizationId: &input.OrganizationId,
		DomainName:     &input.DomainName,
	})
	if err != nil {
		return "", state, err
	}

	mailDomain, err := workmailclient.GetMailDomain(ctx, &workmail.GetMailDomainInput{
		OrganizationId: &input.OrganizationId,
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

	return state.DomainName, state, nil
}

// The Delete method will run when the resource is deleted.
func (DefaultDomain) Delete(ctx p.Context, id string, props DefaultDomainState) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	cfg.Region = props.Region

	// Create the WorkMail service client using the config
	workmailclient := workmail.NewFromConfig(cfg)

	// Reset the default domain
	organization, err := workmailclient.DescribeOrganization(ctx, &workmail.DescribeOrganizationInput{
		OrganizationId: &props.OrganizationId,
	})
	if err != nil {
		return err
	}
	defaultDomain := *organization.Alias + ".awsapps.com"
	_, err = workmailclient.UpdateDefaultMailDomain(ctx, &workmail.UpdateDefaultMailDomainInput{
		OrganizationId: &props.OrganizationId,
		DomainName:     &defaultDomain,
	})
	if err != nil {
		return err
	}

	_, err = workmailclient.DeregisterMailDomain(ctx, &workmail.DeregisterMailDomainInput{
		OrganizationId: &props.OrganizationId,
		DomainName:     &props.DomainName,
	})
	if err != nil {
		return err
	}

	return err
}
