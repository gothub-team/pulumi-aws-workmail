using System.Collections.Generic;
using System.Linq;
using Pulumi;
using awsworkmail = Pulumi.awsworkmail;

return await Deployment.RunAsync(() => 
{
    var myRandomResource = new awsworkmail.Random("myRandomResource", new()
    {
        Length = 24,
    });

    return new Dictionary<string, object?>
    {
        ["output"] = 
        {
            { "value", myRandomResource.Result },
        },
    };
});

