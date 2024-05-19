using System.Collections.Generic;
using System.Linq;
using Pulumi;
using aws-workmail = Pulumi.aws-workmail;

return await Deployment.RunAsync(() => 
{
    var myRandomResource = new aws-workmail.Random("myRandomResource", new()
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

