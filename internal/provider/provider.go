package provider

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				"remotefile": dataSourceRemotefile(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"remotefile": resourceRemotefile(),
			},
			Schema: map[string]*schema.Schema{
				"max_sessions": {
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     3,
					Description: "Maximum number of open sessions in each host connection.",
				},
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

type apiClient struct {
	mux            *sync.Mutex
	remoteClients  map[string]*RemoteClient
	activeSessions map[string]int
	maxSessions    int
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(c context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		client := apiClient{
			maxSessions:    d.Get("max_sessions").(int),
			mux:            &sync.Mutex{},
			remoteClients:  map[string]*RemoteClient{},
			activeSessions: map[string]int{},
		}

		return &client, diag.Diagnostics{}
	}
}

func (c *apiClient) getRemoteClient(d *schema.ResourceData) (*RemoteClient, error) {
	connectionID := resourceConnectionHash(d)
	defer c.mux.Unlock()
	for {
		c.mux.Lock()

		client, ok := c.remoteClients[connectionID]
		if ok {
			if c.activeSessions[connectionID] >= c.maxSessions {
				c.mux.Unlock()
				continue
			}
			c.activeSessions[connectionID] += 1

			return client, nil
		}

		client, err := RemoteClientFromResource(d)
		if err != nil {
			return nil, err
		}

		c.remoteClients[connectionID] = client
		c.activeSessions[connectionID] = 1
		return client, nil
	}
}

func (c *apiClient) closeRemoteClient(d *schema.ResourceData) error {
	connectionID := resourceConnectionHash(d)
	c.mux.Lock()
	defer c.mux.Unlock()

	c.activeSessions[connectionID] -= 1
	if c.activeSessions[connectionID] == 0 {
		client := c.remoteClients[connectionID]
		delete(c.remoteClients, connectionID)
		return client.Close()
	}

	return nil
}

func RemoteClientFromResource(d *schema.ResourceData) (*RemoteClient, error) {
	clientConfig := ssh.ClientConfig{
		User:            d.Get("conn.0.username").(string),
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
			return nil, fmt.Errorf("couldn't create a ssh client config from private key: %s", err.Error())
		}
		clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(signer))
	}

	private_key_path, ok := d.GetOk("conn.0.private_key_path")
	if ok {
		content, err := ioutil.ReadFile(private_key_path.(string))
		if err != nil {
			return nil, fmt.Errorf("couldn't read private key: %s", err.Error())
		}
		signer, err := ssh.ParsePrivateKey(content)
		if err != nil {
			return nil, fmt.Errorf("couldn't create a ssh client config from private key file: %s", err.Error())
		}
		clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(signer))
	}

	private_key_env_var, ok := d.GetOk("conn.0.private_key_env_var")
	if ok {
		private_key := os.Getenv(private_key_env_var.(string))
		signer, err := ssh.ParsePrivateKey([]byte(private_key))
		if err != nil {
			return nil, fmt.Errorf("couldn't create a ssh client config from private key env var: %s", err.Error())
		}
		clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(signer))
	}

	host := fmt.Sprintf("%s:%d", d.Get("conn.0.host").(string), d.Get("conn.0.port").(int))
	return NewRemoteClient(host, clientConfig)
}

func resourceConnectionHash(d *schema.ResourceData) string {
	elements := []string{
		d.Get("conn.0.host").(string),
		d.Get("conn.0.username").(string),
		strconv.Itoa(d.Get("conn.0.port").(int)),
		resourceStringWithDefault(d, "conn.0.password", ""),
		resourceStringWithDefault(d, "conn.0.private_key", ""),
		resourceStringWithDefault(d, "conn.0.private_key_path", ""),
	}
	return strings.Join(elements, "::")
}

func resourceStringWithDefault(d *schema.ResourceData, key string, defaultValue string) string {
	str, ok := d.GetOk(key)
	if ok {
		return str.(string)
	}
	return defaultValue
}
