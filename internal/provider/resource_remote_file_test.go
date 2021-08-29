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

func TestAccResourceRemoteFileWithDefaultConnection(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRemoteFileWithDefaultConnection,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"remote_file.bar", "content", regexp.MustCompile("123")),
				),
			},
		},
	})
}

const testAccResourceRemoteFile = `
resource "remote_file" "foo" {
  conn {
	  host = "remotehost"
	  user = "root"
	  sudo = true
	  password = "password"
  }
  path = "/tmp/foo.txt"
  content = "bar"
  permissions = "0777"
}
`

const testAccResourceRemoteFileWithDefaultConnection = `
resource "remote_file" "bar" {
	provider = remotehost

	path = "/tmp/defaultconn.txt"
	content = "123"
  }
`
