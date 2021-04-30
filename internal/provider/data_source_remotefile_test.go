package provider

import (
	"io/ioutil"
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
				PreConfig: func() {
					err := ioutil.WriteFile("/tmp/bar.txt", []byte("baz"), 0777)
					if err != nil {
						t.Fatal("unable to create test file to read")
					}
				},
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
	path = "/tmp/bar.txt"
}
`
