package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceRemotefile(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRemotefile,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"remotefile_resource.foo", "content", regexp.MustCompile("bar")),
				),
			},
		},
	})
}

const testAccResourceRemotefile = `
resource "remotefile_resource" "foo" {
  path = "/tmp/foo.txt"
  content = "bar"
  permissions = "0777"
}
`
