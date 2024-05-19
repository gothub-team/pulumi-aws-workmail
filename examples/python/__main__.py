import pulumi
import pulumi_aws-workmail as aws-workmail

my_random_resource = aws-workmail.Random("myRandomResource", length=24)
pulumi.export("output", {
    "value": my_random_resource.result,
})
