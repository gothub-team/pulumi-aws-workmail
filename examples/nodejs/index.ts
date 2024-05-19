import * as pulumi from "@pulumi/pulumi";
import * as awsworkmail from "@pulumi/awsworkmail";

const myRandomResource = new awsworkmail.Random("myRandomResource", {length: 24});
export const output = {
    value: myRandomResource.result,
};
