package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"golang.org/x/crypto/ssh"
)

func TestAccDataSourceRemoteFile(t *testing.T) {

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			sshClient, err := ssh.Dial("tcp", "remotehost:22", &ssh.ClientConfig{
				User:            "root",
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
				Auth:            []ssh.AuthMethod{ssh.Password("password")},
			})
			if err != nil {
				panic(err)
			}
			session, err := sshClient.NewSession()
			if err != nil {
				panic(err)
			}

			defer session.Close()

			stdin, err := session.StdinPipe()
			if err != nil {
				panic(err)
			}
			go func() {
				stdin.Write([]byte("file-content"))
				stdin.Close()
			}()

			session.Run(fmt.Sprintf("cat /dev/stdin | tee %s", "/tmp/bar.txt"))
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
