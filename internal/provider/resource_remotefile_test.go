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
				Config: `
				resource "remotefile" "foo" {
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
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"remotefile.foo", "content", regexp.MustCompile("bar")),
				),
			},
			{
				Config: `
				resource "remotefile" "bar" {
				  conn {
					  host = "remotehost8022"
					  port = 8022
					  private_key_path = "../../tests/key"
					  username = "root"
					  sudo = true
				  }
				  path = "/tmp/bar.txt"
				  content = "bar"
				  permissions = "0777"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"remotefile.bar", "content", regexp.MustCompile("bar")),
				),
			},
		},
	})
}
