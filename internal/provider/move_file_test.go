package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Updating fields in provider or swapping provider should ideally be supported
/*func TestMovingFileByModifyingProvider(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				// Create file on 'remotehost'
				Config: `
				resource "remote_file" "move_file" {
					provider = remotehost
					path = "/tmp/move_file.txt"
					content = "x"
				}
				`,
			},
			{
				// Move file to 'remotehost2'
				Config: `
				resource "remote_file" "move_file" {
					provider = remotehost2
					path = "/tmp/move_file.txt"
					content = "x"
				}
				`,
			},
			{
				// Read file on 'remotehost2'
				Config: `
				resource "remote_file" "move_file" {
					provider = remotehost2
					path = "/tmp/move_file.txt"
					content = "x"
				}
				data "remote_file" "move_file" {
					provider = remotehost2
					path = "/tmp/move_file.txt"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.remote_file.move_file", "content", regexp.MustCompile("x")),
				),
			},
		},
	})
}*/

func TestMovingFileByModifyingConn(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				// Create file on 'remotehost'
				Config: `
				resource "remote_file" "move_file_2" {
					conn {
						host = "remotehost"
						user = "root"
						password = "password"
					}
					path = "/tmp/move_file_2.txt"
					content = "x"
				}
				`,
			},
			{
				// Move file to 'remotehost2'
				Config: `
				resource "remote_file" "move_file_2" {
					conn {
						host = "remotehost2"
						user = "root"
						password = "password"
					}
					path = "/tmp/move_file_2.txt"
					content = "x"
				}
				`,
			},
			{
				// Read file on 'remotehost2'
				Config: `
				resource "remote_file" "move_file_2" {
					conn {
						host = "remotehost2"
						user = "root"
						password = "password"
					}
					path = "/tmp/move_file_2.txt"
					content = "x"
				}
				data "remote_file" "move_file_2" {
					conn {
						host = "remotehost2"
						user = "root"
						password = "password"
					}
					path = "/tmp/move_file_2.txt"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.remote_file.move_file_2", "content", regexp.MustCompile("x")),
				),
			},
		},
	})
}
