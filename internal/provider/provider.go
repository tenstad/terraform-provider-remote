package provider

import (
	"context"
	"fmt"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			DataSourcesMap: map[string]*schema.Resource{
				"remotefile_data_source": dataSourceRemotefile(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"remotefile_resource": resourceRemotefile(),
			},
			Schema: map[string]*schema.Schema{
				"username": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("REMOTEFILE_USERNAME", nil),
					Description: "The username on the target host. May alternatively be set via the `REMOTEFILE_USERNAME` environment variable.",
				},
				"private_key_path": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("REMOTEFILE_PRIVATE_KEY_PATH", nil),
					Description: "The path to the private key used to login to target host. May alternatively be set via the `REMOTEFILE_PRIVATE_KEY_PATH` environment variable.",
				},
				"host": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("REMOTEFILE_HOST", nil),
					Description: "The target host where files are located. May alternatively be set via the `REMOTEFILE_HOST` environment variable.",
				},
				"port": {
					Type:        schema.TypeInt,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("REMOTEFILE_PORT", 22),
					Description: "The ssh port to the target host. May alternatively be set via the `REMOTEFILE_PORT` environment variable.",
				},
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

type apiClient struct {
	clientConfig ssh.ClientConfig
	host         string
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(c context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		clientConfig, err := auth.PrivateKey(d.Get("username").(string), d.Get("private_key_path").(string), ssh.InsecureIgnoreHostKey())
		if err != nil {
			return nil, diag.Errorf("couldn't create a ssh client config: %s", err.Error())
		}

		client := apiClient{
			clientConfig: clientConfig,
			host:         fmt.Sprintf("%s:%d", d.Get("host").(string), d.Get("port").(int)),
		}

		return &client, nil
	}
}

func (c apiClient) getSSHClient() (*ssh.Client, error) {
	sshClient, err := ssh.Dial("tcp", c.host, &c.clientConfig)
	if err != nil {
		return nil, fmt.Errorf("couldn't establish a connection to the remote server: %s", err.Error())
	}
	return sshClient, nil
}

func (c apiClient) getSCPClient() (*scp.Client, error) {
	scpClient := scp.NewClient(c.host, &c.clientConfig)
	err := scpClient.Connect()
	if err != nil {
		return nil, fmt.Errorf("couldn't establish a connection to the remote server: %s", err.Error())
	}
	return &scpClient, nil
}

func (c apiClient) getSFTPClient() (*sftp.Client, error) {
	sshClient, err := c.getSSHClient()
	if err != nil {
		return nil, err
	}

	sftp, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, err
	}
	return sftp, nil
}
