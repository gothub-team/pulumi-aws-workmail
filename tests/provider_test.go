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
	. "github.com/smartystreets/goconvey/convey"

	awsworkmail "github.com/gothub-team/pulumi-awsworkmail/provider"
)

func SkipTestRandomCreate(t *testing.T) {
	prov := provider()
	Convey("Given a length", t, func() {
		length := 12

		Convey("When creating a random string", func() {
			response, err := prov.Create(p.CreateRequest{
				Urn: urn("Random"),
				Properties: resource.PropertyMap{
					"length": resource.NewNumberProperty(float64(length)),
				},
				Preview: false,
			})

			So(err, ShouldBeNil)

			result := response.Properties["result"].StringValue()
			So(result, ShouldHaveLength, 12)
		})
	})
}

func TestOrganization(t *testing.T) {
	prov := provider()

	Convey("When creating an organization", t, func() {
		organization, err := prov.Create(p.CreateRequest{
			Urn: urn("Organization"),
			Properties: resource.PropertyMap{
				"region": resource.NewStringProperty("eu-west-1"),
				"alias":  resource.NewStringProperty("test-organization-alias"),
			},
			Preview: false,
		})

		So(err, ShouldBeNil)
		So(organization.Properties["organizationId"].StringValue(), ShouldNotBeEmpty)

		Convey("When creating a default domain", func() {
			domain, err := prov.Create(p.CreateRequest{
				Urn: urn("DefaultDomain"),
				Properties: resource.PropertyMap{
					"region":         resource.NewStringProperty("eu-west-1"),
					"domainName":     resource.NewStringProperty("dev.gothub.io"),
					"organizationId": organization.Properties["organizationId"],
				},
				Preview: false,
			})

			So(err, ShouldBeNil)
			So(domain.Properties["records"].ArrayValue(), ShouldHaveLength, 8)

			SkipConvey("When deleting the default domain", func() {
				err := prov.Delete(p.DeleteRequest{
					Urn:        urn("DefaultDomain"),
					Properties: domain.Properties,
					ID:         domain.ID,
				})

				So(err, ShouldBeNil)
			})
		})

		err = prov.Delete(p.DeleteRequest{
			Urn:        urn("Organization"),
			Properties: organization.Properties,
			ID:         organization.ID,
		})

		So(err, ShouldBeNil)
	})
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
