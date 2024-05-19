package main

import (
	"github.com/gothub-team/pulumi-aws-workmail/sdk/go/aws-workmail"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		myRandomResource, err := aws-workmail.NewRandom(ctx, "myRandomResource", &aws-workmail.RandomArgs{
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
