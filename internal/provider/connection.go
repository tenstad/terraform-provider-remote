package provider

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

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
			ForceNew:    true,
			Description: "The ssh port on the remote host.",
		},
		"timeout": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The maximum amount of time, in milliseconds, for the TCP connection to establish. Timeout of zero means no timeout.",
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
		"agent": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Use a local SSH agent to login to the remote host.",
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
	},
}

func ConnectionFromResourceData(ctx context.Context, d *schema.ResourceData) (string, *ssh.ClientConfig, error) {
	if _, ok := d.GetOk("conn"); !ok {
		return "", nil, fmt.Errorf("resouce does not have a connection configured")
	}

	clientConfig := ssh.ClientConfig{
		User:            d.Get("conn.0.user").(string),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	password, ok := d.GetOk("conn.0.password")
	if ok {
		clientConfig.Auth = append(clientConfig.Auth, ssh.Password(password.(string)))
	}

	private_key, ok := d.GetOk("conn.0.private_key")
	if ok {
		signer, err := ssh.ParsePrivateKey([]byte(private_key.(string)))
		if err != nil {
			return "", nil, fmt.Errorf("couldn't create a ssh client config from private key: %s", err.Error())
		}
		clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(signer))
	}

	private_key_path, ok := d.GetOk("conn.0.private_key_path")
	if ok {
		content, err := os.ReadFile(private_key_path.(string))
		if err != nil {
			return "", nil, fmt.Errorf("couldn't read private key: %s", err.Error())
		}
		signer, err := ssh.ParsePrivateKey(content)
		if err != nil {
			return "", nil, fmt.Errorf("couldn't create a ssh client config from private key file: %s", err.Error())
		}
		clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(signer))
	}

	private_key_env_var, ok := d.GetOk("conn.0.private_key_env_var")
	if ok {
		private_key := os.Getenv(private_key_env_var.(string))
		signer, err := ssh.ParsePrivateKey([]byte(private_key))
		if err != nil {
			return "", nil, fmt.Errorf("couldn't create a ssh client config from private key env var: %s", err.Error())
		}
		clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(signer))
	}

	enableAgent, ok := d.GetOk("conn.0.agent")
	if ok && enableAgent.(bool) {
		connection, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
		if err != nil {
			return "", nil, fmt.Errorf("couldn't connect to SSH agent: %s", err.Error())
		}
		clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeysCallback(agent.NewClient(connection).Signers))
	}

	timeout, ok := d.GetOk("conn.0.timeout")
	if ok {
		clientConfig.Timeout = time.Duration(timeout.(int)) * time.Millisecond
	}

	host := fmt.Sprintf("%s:%d", d.Get("conn.0.host").(string), d.Get("conn.0.port").(int))
	return host, &clientConfig, nil
}
