package provider

import (
	"os"
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
				Config: `
				resource "remote_file" "resource_1" {
				  conn {
					  host = "remotehost"
					  user = "root"
					  sudo = true
					  password = "password"
				  }
				  path = "/tmp/resource_1.txt"
				  content = "resource_1"
				  permissions = "0777"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"remote_file.resource_1", "content", regexp.MustCompile("resource_1")),
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
				Config: `
				resource "remote_file" "resource_2" {
					provider = remotehost
				
					path = "/tmp/resource_2.txt"
					content = "resource_2"
				  }
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"remote_file.resource_2", "id", regexp.MustCompile("remotehost:22:/tmp/resource_2.txt")),
					resource.TestMatchResourceAttr(
						"remote_file.resource_2", "content", regexp.MustCompile("resource_2")),
				),
			},
		},
	})
}

func TestAccResourceRemoteFileWithAgent(t *testing.T) {
	if os.Getenv("SKIP_TEST_AGENT") == "1" {
		return
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "remote_file" "resource_3" {
					conn {
						host = "remotehost"
						user = "root"
						agent = true
					}
					path = "/tmp/resource_3.txt"
					content = "resource_3"
				  }
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"remote_file.resource_3", "content", regexp.MustCompile("resource_3")),
				),
			},
		},
	})
}
