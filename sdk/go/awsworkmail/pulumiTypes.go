// Code generated by pulumi-language-go DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package awsworkmail

import (
	"context"
	"reflect"

	"github.com/gothub-team/pulumi-awsworkmail/sdk/go/awsworkmail/internal"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumix"
)

var _ = internal.GetEnvOrDefault

type DomainArgs struct {
	DomainName   string `pulumi:"domainName"`
	HostedZoneId string `pulumi:"hostedZoneId"`
}

type DomainArgsArgs struct {
	DomainName   pulumix.Input[string] `pulumi:"domainName"`
	HostedZoneId pulumix.Input[string] `pulumi:"hostedZoneId"`
}

func (DomainArgsArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*DomainArgs)(nil)).Elem()
}

func (i DomainArgsArgs) ToDomainArgsOutput() DomainArgsOutput {
	return i.ToDomainArgsOutputWithContext(context.Background())
}

func (i DomainArgsArgs) ToDomainArgsOutputWithContext(ctx context.Context) DomainArgsOutput {
	return pulumi.ToOutputWithContext(ctx, i).(DomainArgsOutput)
}

func (i *DomainArgsArgs) ToOutput(ctx context.Context) pulumix.Output[*DomainArgsArgs] {
	return pulumix.Val(i)
}

type DomainArgsOutput struct{ *pulumi.OutputState }

func (DomainArgsOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*DomainArgs)(nil)).Elem()
}

func (o DomainArgsOutput) ToDomainArgsOutput() DomainArgsOutput {
	return o
}

func (o DomainArgsOutput) ToDomainArgsOutputWithContext(ctx context.Context) DomainArgsOutput {
	return o
}

func (o DomainArgsOutput) ToOutput(ctx context.Context) pulumix.Output[DomainArgs] {
	return pulumix.Output[DomainArgs]{
		OutputState: o.OutputState,
	}
}

func (o DomainArgsOutput) DomainName() pulumix.Output[string] {
	return pulumix.Apply[DomainArgs](o, func(v DomainArgs) string { return v.DomainName })
}

func (o DomainArgsOutput) HostedZoneId() pulumix.Output[string] {
	return pulumix.Apply[DomainArgs](o, func(v DomainArgs) string { return v.HostedZoneId })
}

func init() {
	pulumi.RegisterOutputType(DomainArgsOutput{})
}
