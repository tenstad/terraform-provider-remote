package provider

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

var connectionSchemaResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"host": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The remote host.",
		},
		"port": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     22,
			Description: "The ssh port on the remote host.",
		},
		"user": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The user on the remote host.",
		},
		"sudo": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Use sudo to gain access to file.",
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "The pasword for the user on the remote host.",
		},
		"private_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "The private key used to login to the remote host.",
		},
		"private_key_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The local path to the private key used to login to the remote host.",
		},
		"private_key_env_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the local environment variable containing the private key used to login to the remote host.",
		},
		"agent": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Use a local SSH agent to login to the remote host.",
		},

	},
}

func ConnectionFromResourceData(d *schema.ResourceData) (string, *ssh.ClientConfig, error) {
	_, ok := d.GetOk("result_conn")
	if !ok {
		return "", nil, fmt.Errorf("resouce does not have a connection configured")
	}

	clientConfig := ssh.ClientConfig{
		User:            d.Get("result_conn.0.user").(string),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	password, ok := d.GetOk("result_conn.0.password")
	if ok {
		clientConfig.Auth = append(clientConfig.Auth, ssh.Password(password.(string)))
	}

	private_key, ok := d.GetOk("result_conn.0.private_key")
	if ok {
		signer, err := ssh.ParsePrivateKey([]byte(private_key.(string)))
		if err != nil {
			return "", nil, fmt.Errorf("couldn't create a ssh client config from private key: %s", err.Error())
		}
		clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(signer))
	}

	private_key_path, ok := d.GetOk("result_conn.0.private_key_path")
	if ok {
		content, err := ioutil.ReadFile(private_key_path.(string))
		if err != nil {
			return "", nil, fmt.Errorf("couldn't read private key: %s", err.Error())
		}
		signer, err := ssh.ParsePrivateKey(content)
		if err != nil {
			return "", nil, fmt.Errorf("couldn't create a ssh client config from private key file: %s", err.Error())
		}
		clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(signer))
	}

	private_key_env_var, ok := d.GetOk("result_conn.0.private_key_env_var")
	if ok {
		private_key := os.Getenv(private_key_env_var.(string))
		signer, err := ssh.ParsePrivateKey([]byte(private_key))
		if err != nil {
			return "", nil, fmt.Errorf("couldn't create a ssh client config from private key env var: %s", err.Error())
		}
		clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(signer))
	}

	enable_agent, ok := d.GetOk("result_conn.0.agent")
	if ok && enable_agent.(bool) {
		connection, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
		if err != nil {
			return "", nil, fmt.Errorf("couldn't connect to SSH agent: %s", err.Error())
		}
		clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeysCallback(agent.NewClient(connection).Signers))
	}

	host := fmt.Sprintf("%s:%d", d.Get("result_conn.0.host").(string), d.Get("result_conn.0.port").(int))
	return host, &clientConfig, nil
}
