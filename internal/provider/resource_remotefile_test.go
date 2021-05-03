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
						"remotefile.foo", "content", regexp.MustCompile("bar")),
				),
			},
		},
	})
}

const testAccResourceRemotefile = `
resource "remotefile" "foo" {
  conn {
	  host = "terraform-provider-remotefile-test"
	  username = "root"
	  sudo = true
	  password = "password"
	  port = 22
  }
  path = "/tmp/foo.txt"
  content = "bar"
  permissions = "0777"
}
`
