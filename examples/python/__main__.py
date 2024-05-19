import pulumi
import pulumi_awsworkmail as awsworkmail

my_random_resource = awsworkmail.Random("myRandomResource", length=24)
pulumi.export("output", {
    "value": my_random_resource.result,
})
