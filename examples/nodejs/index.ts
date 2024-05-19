import * as pulumi from "@pulumi/pulumi";
import * as aws-workmail from "@pulumi/aws-workmail";

const myRandomResource = new aws-workmail.Random("myRandomResource", {length: 24});
export const output = {
    value: myRandomResource.result,
};
