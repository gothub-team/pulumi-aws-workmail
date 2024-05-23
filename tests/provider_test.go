// Copyright 2016-2023, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tests

import (
	"testing"

	"github.com/blang/semver"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	awsworkmail "github.com/gothub-team/pulumi-awsworkmail/provider"
)

func TestRandomCreate(t *testing.T) {
	prov := provider()

	response, err := prov.Create(p.CreateRequest{
		Urn: urn("Random"),
		Properties: resource.PropertyMap{
			"length": resource.NewNumberProperty(12),
		},
		Preview: false,
	})

	require.NoError(t, err)
	result := response.Properties["result"].StringValue()
	assert.Len(t, result, 12)
}

func TestOrganizationCreateDelete(t *testing.T) {
	prov := provider()

	err := prov.Delete(p.DeleteRequest{
		Urn: urn("Organization"),
		Properties: resource.PropertyMap{
			"organizationId": resource.NewStringProperty("m-9300c47ab20d4182ab1c4c3ebb8d358b"),
		}})
	require.NoError(t, err)

	// response, err := prov.Create(p.CreateRequest{
	// 	Urn: urn("Organization"),
	// 	Properties: resource.PropertyMap{
	// 		"region":       resource.NewStringProperty("eu-west-1"),
	// 		"alias":        resource.NewStringProperty("test-devgothubio"),
	// 		"domainName":   resource.NewStringProperty("dev.gothub.io"),
	// 		"hostedZoneId": resource.NewStringProperty("Z0690737HWV9262JDHN4"),
	// 	},
	// 	Preview: false,
	// })

	// response.Properties["organizationId"].StringValue()

	// require.NoError(t, err)
	// result := response.Properties["organizationId"].StringValue()
	// assert.NotEmpty(t, result)
}

// urn is a helper function to build an urn for running integration tests.
func urn(typ string) resource.URN {
	return resource.NewURN("stack", "proj", "",
		tokens.Type("test:index:"+typ), "name")
}

// Create a test server.
func provider() integration.Server {
	return integration.NewServer(awsworkmail.Name, semver.MustParse("1.0.0"), awsworkmail.Provider())
}
