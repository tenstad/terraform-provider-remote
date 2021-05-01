package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRemotefile(t *testing.T) {

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRemotefile,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						// TODO: check content is correct
						"data.remotefile.bar", "content", regexp.MustCompile("")),
				),
			},
		},
	})
}

const testAccDataSourceRemotefile = `
data "remotefile" "bar" {
	conn {
		host = "localhost"
		username = "root"
		private_key_path = "../../tests/key"
		port = 8022
	}
	path = "/tmp/bar.txt"
}
`
