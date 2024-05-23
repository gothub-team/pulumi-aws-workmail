package provider

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/workmail"
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
type WorkmailRegistration struct{}

// Each resource has an input struct, defining what arguments it accepts.
type WorkmailRegistrationArgs struct {
	// The AWS Region. TODO: This should be passed as a pulumi.Provider
	Region string `pulumi:"region"`
	// The organization id.
	OrganizationId string `pulumi:"organizationId"`
	// The identifier for the user, group, or resource to be updated.
	// The identifier can accept UserId, ResourceId, or GroupId, or Username, Resourcename, or Groupname. The following identity formats are available:
	// Entity ID: 12345678-1234-1234-1234-123456789012, r-0123456789a0123456789b0123456789, or S-1-1-12-1234567890-123456789-123456789-1234
	// Entity name: entity
	// This member is required.
	EntityId string `pulumi:"entityId"`
	// The email prefix for the new user. (prefix@domain.com).
	// The default domain of the organization will be appended automatically.
	EmailPrefix string `pulumi:"emailPrefix"`
}

// Each resource has a state, describing the fields that exist on the created resource.
type WorkmailRegistrationState struct {
	// It is generally a good idea to embed args in outputs, but it isn't strictly necessary.
	WorkmailRegistrationArgs
}

// All resources must implement Create at a minimum.
func (WorkmailRegistration) Create(ctx p.Context, name string, input WorkmailRegistrationArgs, preview bool) (string, WorkmailRegistrationState, error) {
	state := WorkmailRegistrationState{WorkmailRegistrationArgs: input}
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

	organization, err := workmailclient.DescribeOrganization(ctx, &workmail.DescribeOrganizationInput{
		OrganizationId: &state.OrganizationId,
	})
	if err != nil {
		return "", state, err
	}

	emailAddress := input.EmailPrefix + "@" + *organization.DefaultMailDomain
	fmt.Println("Creating primary email address", emailAddress)
	_, err = workmailclient.RegisterToWorkMail(ctx, &workmail.RegisterToWorkMailInput{
		OrganizationId: &state.OrganizationId,
		EntityId:       &input.EntityId,
		Email:          &emailAddress,
	})
	if err != nil {
		return "", state, err
	}

	return input.EntityId, state, nil
}

// The Delete method will run when the resource is deleted.
func (WorkmailRegistration) Delete(ctx p.Context, id string, props WorkmailRegistrationState) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	cfg.Region = props.Region

	// Create the WorkMail service client using the config
	workmailclient := workmail.NewFromConfig(cfg)

	_, err = workmailclient.DeregisterFromWorkMail(ctx, &workmail.DeregisterFromWorkMailInput{
		OrganizationId: &props.OrganizationId,
		EntityId:       &props.EntityId,
	})

	return err
}
