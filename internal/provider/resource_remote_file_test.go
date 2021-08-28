package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceRemoteFile(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRemoteFile,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"remote_file.foo", "content", regexp.MustCompile("bar")),
				),
			},
		},
	})
}

const testAccResourceRemoteFile = `
resource "remote_file" "foo" {
  conn {
	  host = "remotehost"
	  username = "root"
	  sudo = true
	  password = "password"
  }
  path = "/tmp/foo.txt"
  content = "bar"
  permissions = "0777"
}
`
