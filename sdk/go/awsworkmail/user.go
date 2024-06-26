// Code generated by pulumi-language-go DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package awsworkmail

import (
	"context"
	"reflect"

	"errors"
	"github.com/gothub-team/pulumi-awsworkmail/sdk/go/awsworkmail/internal"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumix"
)

type User struct {
	pulumi.CustomResourceState

	DisplayName                 pulumix.Output[string]  `pulumi:"displayName"`
	Domain                      pulumix.Output[*string] `pulumi:"domain"`
	FirstName                   pulumix.Output[*string] `pulumi:"firstName"`
	HiddenFromGlobalAddressList pulumix.Output[*bool]   `pulumi:"hiddenFromGlobalAddressList"`
	LastName                    pulumix.Output[*string] `pulumi:"lastName"`
	Name                        pulumix.Output[string]  `pulumi:"name"`
	OrganizationId              pulumix.Output[string]  `pulumi:"organizationId"`
	Password                    pulumix.Output[*string] `pulumi:"password"`
	Region                      pulumix.Output[string]  `pulumi:"region"`
	UserId                      pulumix.Output[string]  `pulumi:"userId"`
}

// NewUser registers a new resource with the given unique name, arguments, and options.
func NewUser(ctx *pulumi.Context,
	name string, args *UserArgs, opts ...pulumi.ResourceOption) (*User, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.DisplayName == nil {
		return nil, errors.New("invalid value for required argument 'DisplayName'")
	}
	if args.Name == nil {
		return nil, errors.New("invalid value for required argument 'Name'")
	}
	if args.Region == nil {
		return nil, errors.New("invalid value for required argument 'Region'")
	}
	opts = internal.PkgResourceDefaultOpts(opts)
	var resource User
	err := ctx.RegisterResource("awsworkmail:index:User", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetUser gets an existing User resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetUser(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *UserState, opts ...pulumi.ResourceOption) (*User, error) {
	var resource User
	err := ctx.ReadResource("awsworkmail:index:User", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering User resources.
type userState struct {
}

type UserState struct {
}

func (UserState) ElementType() reflect.Type {
	return reflect.TypeOf((*userState)(nil)).Elem()
}

type userArgs struct {
	DisplayName                 string  `pulumi:"displayName"`
	Domain                      *string `pulumi:"domain"`
	FirstName                   *string `pulumi:"firstName"`
	HiddenFromGlobalAddressList *bool   `pulumi:"hiddenFromGlobalAddressList"`
	LastName                    *string `pulumi:"lastName"`
	Name                        string  `pulumi:"name"`
	OrganizationId              *string `pulumi:"organizationId"`
	Password                    *string `pulumi:"password"`
	Region                      string  `pulumi:"region"`
}

// The set of arguments for constructing a User resource.
type UserArgs struct {
	DisplayName                 pulumix.Input[string]
	Domain                      pulumix.Input[*string]
	FirstName                   pulumix.Input[*string]
	HiddenFromGlobalAddressList pulumix.Input[*bool]
	LastName                    pulumix.Input[*string]
	Name                        pulumix.Input[string]
	OrganizationId              pulumix.Input[*string]
	Password                    pulumix.Input[*string]
	Region                      pulumix.Input[string]
}

func (UserArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*userArgs)(nil)).Elem()
}

type UserOutput struct{ *pulumi.OutputState }

func (UserOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*User)(nil)).Elem()
}

func (o UserOutput) ToUserOutput() UserOutput {
	return o
}

func (o UserOutput) ToUserOutputWithContext(ctx context.Context) UserOutput {
	return o
}

func (o UserOutput) ToOutput(ctx context.Context) pulumix.Output[User] {
	return pulumix.Output[User]{
		OutputState: o.OutputState,
	}
}

func (o UserOutput) DisplayName() pulumix.Output[string] {
	value := pulumix.Apply[User](o, func(v User) pulumix.Output[string] { return v.DisplayName })
	return pulumix.Flatten[string, pulumix.Output[string]](value)
}

func (o UserOutput) Domain() pulumix.Output[*string] {
	value := pulumix.Apply[User](o, func(v User) pulumix.Output[*string] { return v.Domain })
	return pulumix.Flatten[*string, pulumix.Output[*string]](value)
}

func (o UserOutput) FirstName() pulumix.Output[*string] {
	value := pulumix.Apply[User](o, func(v User) pulumix.Output[*string] { return v.FirstName })
	return pulumix.Flatten[*string, pulumix.Output[*string]](value)
}

func (o UserOutput) HiddenFromGlobalAddressList() pulumix.Output[*bool] {
	value := pulumix.Apply[User](o, func(v User) pulumix.Output[*bool] { return v.HiddenFromGlobalAddressList })
	return pulumix.Flatten[*bool, pulumix.Output[*bool]](value)
}

func (o UserOutput) LastName() pulumix.Output[*string] {
	value := pulumix.Apply[User](o, func(v User) pulumix.Output[*string] { return v.LastName })
	return pulumix.Flatten[*string, pulumix.Output[*string]](value)
}

func (o UserOutput) Name() pulumix.Output[string] {
	value := pulumix.Apply[User](o, func(v User) pulumix.Output[string] { return v.Name })
	return pulumix.Flatten[string, pulumix.Output[string]](value)
}

func (o UserOutput) OrganizationId() pulumix.Output[string] {
	value := pulumix.Apply[User](o, func(v User) pulumix.Output[string] { return v.OrganizationId })
	return pulumix.Flatten[string, pulumix.Output[string]](value)
}

func (o UserOutput) Password() pulumix.Output[*string] {
	value := pulumix.Apply[User](o, func(v User) pulumix.Output[*string] { return v.Password })
	return pulumix.Flatten[*string, pulumix.Output[*string]](value)
}

func (o UserOutput) Region() pulumix.Output[string] {
	value := pulumix.Apply[User](o, func(v User) pulumix.Output[string] { return v.Region })
	return pulumix.Flatten[string, pulumix.Output[string]](value)
}

func (o UserOutput) UserId() pulumix.Output[string] {
	value := pulumix.Apply[User](o, func(v User) pulumix.Output[string] { return v.UserId })
	return pulumix.Flatten[string, pulumix.Output[string]](value)
}

func init() {
	pulumi.RegisterOutputType(UserOutput{})
}
