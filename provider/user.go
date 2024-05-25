package provider

import (
	"errors"
	"fmt"

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
type User struct{}

// Each resource has an input struct, defining what arguments it accepts.
type UserArgs struct {
	// The AWS Region. TODO: This should be passed as a pulumi.Provider
	Region string `pulumi:"region"`
	// The display name for the new user.
	DisplayName string `pulumi:"displayName"`
	// The name for the new user. WorkMail directory user names have a maximum length
	// of 64. All others have a maximum length of 20.
	Name string `pulumi:"name"`
	// The identifier of the organization for which the user is created. Either
	// organizationId or domain must be specified. When both are specified, organizationId
	// is ignored.
	OrganizationId *string `pulumi:"organizationId,optional"`
	// The mail domain of the organization for which the user is created. Either
	// organizationId or domain must be specified. When both are specified, organizationId
	// is ignored.
	Domain *string `pulumi:"domain,optional"`
	// The first name of the new user.
	FirstName *string `pulumi:"firstName,optional"`
	// If this parameter is enabled, the user will be hidden from the address book.
	HiddenFromGlobalAddressList *bool `pulumi:"hiddenFromGlobalAddressList,optional"`
	// The last name of the new user.
	LastName *string `pulumi:"lastName,optional"`
	// The password for the new user.
	Password *string `pulumi:"password,optional"`
	// The role of the new user.
	//
	// You cannot pass SYSTEM_USER or RESOURCE role in a single request. When a user
	// role is not selected, the default role of USER is selected.
	// Role *UserRole `pulumi:"role,optional"`
	// contains filtered or unexported fields
}

// type UserRole string

// // Enum values for UserRole
// const (
// 	UserRoleUser       UserRole = "USER"
// 	UserRoleResource   UserRole = "RESOURCE"
// 	UserRoleSystemUser UserRole = "SYSTEM_USER"
// 	UserRoleRemoteUser UserRole = "REMOTE_USER"
// )

// Each resource has a state, describing the fields that exist on the created resource.
type UserState struct {
	// It is generally a good idea to embed args in outputs, but it isn't strictly necessary.
	UserArgs

	// The user id.
	UserId string `pulumi:"userId"`
	// The organization id.
	OrganizationId string `pulumi:"organizationId"`
}

// All resources must implement Create at a minimum.
func (User) Create(ctx p.Context, name string, input UserArgs, preview bool) (string, UserState, error) {
	state := UserState{UserArgs: input}
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

	// Find organization
	if input.OrganizationId == nil && input.Domain == nil {
		return "", state, errors.New("either organizationId or domain must be specified")
	}
	if input.Domain != nil {
		organizations, err := workmailclient.ListOrganizations(ctx, &workmail.ListOrganizationsInput{})
		if err != nil {
			return "", state, err
		}
		organization, found := Find(func(org types.OrganizationSummary) bool {
			if org.DefaultMailDomain == nil || input.Domain == nil {
				return false
			}
			return *org.DefaultMailDomain == *input.Domain
		})(organizations.OrganizationSummaries)
		if !found {
			return "", state, fmt.Errorf("no workmail organization with domain %s found", *input.Domain)
		}
		state.OrganizationId = *organization.OrganizationId
	} else {
		state.OrganizationId = *input.OrganizationId
	}

	// Create the organization
	user, err := workmailclient.CreateUser(ctx, &workmail.CreateUserInput{
		DisplayName:                 &input.DisplayName,
		OrganizationId:              &state.OrganizationId,
		Name:                        &input.Name,
		FirstName:                   input.FirstName,
		LastName:                    input.LastName,
		Password:                    input.Password,
		HiddenFromGlobalAddressList: ifNotNil(input.HiddenFromGlobalAddressList, false),
		// Role:                        types.UserRoleUser,
	})
	if err != nil {
		return "", state, err
	}

	state.UserId = *user.UserId

	return *user.UserId, state, nil
}

func (User) Diff(ctx p.Context, id string, olds UserState, news UserArgs) (p.DiffResponse, error) {
	diffs := make(map[string]p.PropertyDiff)
	hasChanges := false

	if news.OrganizationId == nil && news.Domain == nil {
		return p.DiffResponse{}, errors.New("either organizationId or domain must be specified")
	}

	// Old has domain
	//   New has domain
	// no change or test equality of domain
	if olds.Domain != nil && news.Domain != nil {
		if *olds.Domain != *news.Domain {
			diffs["domain"] = p.PropertyDiff{Kind: p.UpdateReplace, InputDiff: true}
			hasChanges = true
		}
	}

	// Old has domain
	//   New has no domain
	// change
	if olds.Domain != nil && news.Domain == nil {
		diffs["domain"] = p.PropertyDiff{Kind: p.DeleteReplace, InputDiff: true}
		diffs["organizationId"] = p.PropertyDiff{Kind: p.AddReplace, InputDiff: true}
		hasChanges = true
	}

	// Old has no domain
	//   New has domain
	// change
	if olds.Domain == nil && news.Domain != nil {
		diffs["domain"] = p.PropertyDiff{Kind: p.AddReplace, InputDiff: true}
		diffs["organizationId"] = p.PropertyDiff{Kind: p.DeleteReplace, InputDiff: true}
		hasChanges = true
	}

	// Old has no domain
	//   New has no domain
	// no change or test equality of organizationId
	if olds.Domain == nil && news.Domain == nil {
		if olds.OrganizationId != *news.OrganizationId {
			diffs["organizationId"] = p.PropertyDiff{Kind: p.UpdateReplace, InputDiff: true}
			hasChanges = true
		}
	}

	//  Region
	if ptrDiff(&olds.Region, &news.Region) {
		diffs["region"] = p.PropertyDiff{Kind: p.UpdateReplace, InputDiff: true}
		hasChanges = true
	}

	//  DisplayName
	if ptrDiff(&olds.DisplayName, &news.DisplayName) {
		diffs["displayName"] = p.PropertyDiff{Kind: p.Update, InputDiff: true}
		hasChanges = true
	}

	//  Name
	if ptrDiff(&olds.Name, &news.Name) {
		diffs["name"] = p.PropertyDiff{Kind: p.Update, InputDiff: true}
		hasChanges = true
	}

	//  FirstName
	if ptrDiff(olds.FirstName, news.FirstName) {
		diffs["firstName"] = p.PropertyDiff{Kind: p.Update, InputDiff: true}
		hasChanges = true
	}

	//  HiddenFromGlobalAddressList
	if ptrDiff(olds.HiddenFromGlobalAddressList, news.HiddenFromGlobalAddressList) {
		diffs["hiddenFromGlobalAddressList"] = p.PropertyDiff{Kind: p.Update, InputDiff: true}
		hasChanges = true
	}

	//  LastName
	if ptrDiff(olds.LastName, news.LastName) {
		diffs["lastName"] = p.PropertyDiff{Kind: p.Update, InputDiff: true}
		hasChanges = true
	}

	//  Password
	if ptrDiff(olds.Password, news.Password) {
		diffs["password"] = p.PropertyDiff{Kind: p.Update, InputDiff: true}
		hasChanges = true
	}

	return p.DiffResponse{HasChanges: hasChanges, DetailedDiff: diffs, DeleteBeforeReplace: true}, nil
}

func ptrDiff[T comparable](a, b *T) bool {
	if a == nil && b == nil {
		return false
	}
	if a == nil || b == nil {
		return true
	}
	return *a != *b
}

// The Delete method will run when the resource is deleted.
func (User) Delete(ctx p.Context, id string, props UserState) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	cfg.Region = props.Region

	// Create the WorkMail service client using the config
	workmailclient := workmail.NewFromConfig(cfg)

	user, err := workmailclient.DescribeUser(ctx, &workmail.DescribeUserInput{
		OrganizationId: &props.OrganizationId,
		UserId:         &props.UserId,
	})
	if err != nil {
		return nil
	}

	if user.State == "DISABLED" {
		_, err = workmailclient.DeleteUser(ctx, &workmail.DeleteUserInput{
			OrganizationId: &props.OrganizationId,
			UserId:         &props.UserId,
		})
	} else {
		return fmt.Errorf("user %s is not in a deletable state (needs to be DISABLED)", props.UserId)
	}
	return err
}

// Find returns a function that takes a slice of type T and returns the first element
// that satisfies the predicate function, along with a boolean indicating if an element was found.
func Find[T any](predicate func(T) bool) func([]T) (T, bool) {
	return func(slice []T) (T, bool) {
		var zero T // zero value for type T
		for _, element := range slice {
			if predicate(element) {
				return element, true
			}
		}
		return zero, false // return zero value and false if no element is found
	}
}
