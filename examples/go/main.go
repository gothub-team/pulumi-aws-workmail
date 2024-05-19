package main

import (
	"github.com/gothub-team/pulumi-awsworkmail/sdk/go/awsworkmail"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		myRandomResource, err := awsworkmail.NewRandom(ctx, "myRandomResource", &awsworkmail.RandomArgs{
			Length: pulumi.Int(24),
		})
		if err != nil {
			return err
		}
		ctx.Export("output", map[string]interface{}{
			"value": myRandomResource.Result,
		})
		return nil
	})
}
