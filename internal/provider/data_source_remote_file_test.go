package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRemoteFile(t *testing.T) {

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			writeFileToHost("remotehost:22", "/tmp/bar.txt", "file-content")
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRemoteFile,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.remote_file.bar", "content", regexp.MustCompile("file-content")),
				),
			},
		},
	})
}

func TestAccDataSourceRemoteFileOverridingDefaultConnection(t *testing.T) {

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			writeFileToHost("remotehost2:22", "/tmp/bar.txt", "file-content")
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRemoteFileOverridingDefaultConnection,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.remote_file.baz", "content", regexp.MustCompile("file-content")),
				),
			},
		},
	})
}

const testAccDataSourceRemoteFile = `
data "remote_file" "bar" {
	conn {
		host = "remotehost"
		user = "root"
		private_key_path = "../../tests/key"
	}
	path = "/tmp/bar.txt"
}
`

const testAccDataSourceRemoteFileOverridingDefaultConnection = `
data "remote_file" "baz" {
	provider = remotehost

	conn {
		host = "remotehost2"
		user = "root"
		private_key_path = "../../tests/key"
	}
	path = "/tmp/bar.txt"
}
`
