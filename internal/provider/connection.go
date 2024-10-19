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

	host, err := Get[string](d, "conn.0.host")
	if err != nil {
		return "", nil, err
	}

	port, err := Get[int](d, "conn.0.port")
	if err != nil {
		return "", nil, err
	}

	user, err := Get[string](d, "conn.0.user")
	if err != nil {
		return "", nil, err
	}

	clientConfig := ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if password, ok, err := GetOk[string](d, "conn.0.password"); ok {
		if err != nil {
			return "", nil, err
		}

		clientConfig.Auth = append(clientConfig.Auth, ssh.Password(password))
	}

	if privateKey, ok, err := GetOk[string](d, "conn.0.private_key"); ok {
		if err != nil {
			return "", nil, err
		}

		signer, err := ssh.ParsePrivateKey([]byte(privateKey))
		if err != nil {
			return "", nil, fmt.Errorf("couldn't create a ssh client config from private key: %s", err.Error())
		}
		clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(signer))
	}

	if privateKeyPath, ok, err := GetOk[string](d, "conn.0.private_key_path"); ok {
		if err != nil {
			return "", nil, err
		}

		content, err := os.ReadFile(privateKeyPath)
		if err != nil {
			return "", nil, fmt.Errorf("couldn't read private key: %s", err.Error())
		}
		signer, err := ssh.ParsePrivateKey(content)
		if err != nil {
			return "", nil, fmt.Errorf("couldn't create a ssh client config from private key file: %s", err.Error())
		}
		clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(signer))
	}

	if privateKeyEnvVar, ok, err := GetOk[string](d, "conn.0.private_key_env_var"); ok {
		if err != nil {
			return "", nil, err
		}

		content := []byte(os.Getenv(privateKeyEnvVar))
		signer, err := ssh.ParsePrivateKey(content)
		if err != nil {
			return "", nil, fmt.Errorf("couldn't create a ssh client config from private key env var: %s", err.Error())
		}
		clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(signer))
	}

	// Don't check ok as terraform struggles with zero values.
	if enableAgent, _, err := GetOk[bool](d, "conn.0.agent"); enableAgent {
		if err != nil {
			return "", nil, err
		}

		connection, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
		if err != nil {
			return "", nil, fmt.Errorf("couldn't connect to SSH agent: %s", err.Error())
		}
		clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeysCallback(agent.NewClient(connection).Signers))
	}

	if timeout, ok, err := GetOk[int](d, "conn.0.timeout"); ok {
		if err != nil {
			return "", nil, err
		}

		clientConfig.Timeout = time.Duration(timeout) * time.Millisecond
	}

	return fmt.Sprintf("%s:%d", host, port), &clientConfig, nil
}
