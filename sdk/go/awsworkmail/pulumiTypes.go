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

type DnsRecord struct {
	Hostname string `pulumi:"hostname"`
	Type     string `pulumi:"type"`
	Value    string `pulumi:"value"`
}

type DnsRecordOutput struct{ *pulumi.OutputState }

func (DnsRecordOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*DnsRecord)(nil)).Elem()
}

func (o DnsRecordOutput) ToDnsRecordOutput() DnsRecordOutput {
	return o
}

func (o DnsRecordOutput) ToDnsRecordOutputWithContext(ctx context.Context) DnsRecordOutput {
	return o
}

func (o DnsRecordOutput) ToOutput(ctx context.Context) pulumix.Output[DnsRecord] {
	return pulumix.Output[DnsRecord]{
		OutputState: o.OutputState,
	}
}

func (o DnsRecordOutput) Hostname() pulumix.Output[string] {
	return pulumix.Apply[DnsRecord](o, func(v DnsRecord) string { return v.Hostname })
}

func (o DnsRecordOutput) Type() pulumix.Output[string] {
	return pulumix.Apply[DnsRecord](o, func(v DnsRecord) string { return v.Type })
}

func (o DnsRecordOutput) Value() pulumix.Output[string] {
	return pulumix.Apply[DnsRecord](o, func(v DnsRecord) string { return v.Value })
}

func init() {
	pulumi.RegisterOutputType(DnsRecordOutput{})
}