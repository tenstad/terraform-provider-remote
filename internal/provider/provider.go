package provider

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure RemoteProvider satisfies various provider interfaces.
var _ provider.Provider = &RemoteProvider{}

// RemoteProvider defines the provider implementation.
type RemoteProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// RemoteProviderModel describes the provider data model.
type RemoteProviderModel struct {
	Conn        types.Int64 `tfsdk:"conn"`
	MaxSessions types.Int64 `tfsdk:"max_sessions"`
}

func (p *RemoteProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "remote"
	resp.Version = p.version
}

func (p *RemoteProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"max_sessions": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of open sessions in each host connection.",
				Optional:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"conn": schema.ListNestedBlock{
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				MarkdownDescription: "Default connection to host where files are located. Can be overridden in resources and data sources.",
				NestedObject:        connectionSchemaResource,
			},
		},
	}
}

func (p *RemoteProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data RemoteProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.MaxSessions.IsNull() {
		data.MaxSessions = basetypes.NewInt64Value(3)
	}

	client := apiClient{
		resourceData:   data,
		maxSessions:    int(data.MaxSessions.ValueInt64()),
		mux:            &sync.Mutex{},
		remoteClients:  map[string]*RemoteClient{},
		activeSessions: map[string]int{},
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *RemoteProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRemoteFileResource,
	}
}

func (p *RemoteProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewRemoteFileDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &RemoteProvider{
			version: version,
		}
	}
}

type apiClient struct {
	resourceData   RemoteProviderModel
	mux            *sync.Mutex
	remoteClients  map[string]*RemoteClient
	activeSessions map[string]int
	maxSessions    int
}

func (c *apiClient) getConnWithDefault(d *schema.ResourceData) (*schema.ResourceData, error) {
	_, ok := d.GetOk("conn")
	if ok {
		return d, nil
	}

	c.mux.Lock()
	defer c.mux.Unlock()

	_, ok = c.resourceData.GetOk("conn")
	if ok {
		return c.resourceData, nil
	}

	return nil, errors.New("neither the provider nor the resource/data source have a configured connection")
}

func (c *apiClient) getRemoteClient(ctx context.Context, d *schema.ResourceData) (*RemoteClient, error) {
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

		client, err := remoteClientFromResourceData(ctx, d)
		if err != nil {
			return nil, err
		}

		c.remoteClients[connectionID] = client
		c.activeSessions[connectionID] = 1
		return client, nil
	}
}

func remoteClientFromResourceData(ctx context.Context, d *schema.ResourceData) (*RemoteClient, error) {
	host, clientConfig, err := ConnectionFromResourceData(ctx, d)
	if err != nil {
		return nil, err
	}
	return NewRemoteClient(host, clientConfig)
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

func setResourceID(d *schema.ResourceData, conn *schema.ResourceData) {
	id := fmt.Sprintf("%s:%d:%s",
		conn.Get("conn.0.host").(string),
		conn.Get("conn.0.port").(int),
		d.Get("path").(string))
	d.SetId(id)
}

func resourceConnectionHash(d *schema.ResourceData) string {
	elements := []string{
		d.Get("conn.0.host").(string),
		d.Get("conn.0.user").(string),
		strconv.Itoa(d.Get("conn.0.port").(int)),
		resourceStringWithDefault(d, "conn.0.password", ""),
		resourceStringWithDefault(d, "conn.0.private_key", ""),
		resourceStringWithDefault(d, "conn.0.private_key_path", ""),
		strconv.FormatBool(d.Get("conn.0.agent").(bool)),
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
