package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRemoteFile(t *testing.T) {

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRemoteFile,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						// TODO: check content is correct
						"data.remote_file.bar", "content", regexp.MustCompile("")),
				),
			},
		},
	})
}

const testAccDataSourceRemoteFile = `
data "remote_file" "bar" {
	conn {
		host = "remotehost"
		username = "root"
		private_key_path = "../../tests/key"
	}
	path = "/tmp/bar.txt"
}
`
