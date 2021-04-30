package provider

import (
	"context"
	"fmt"
	"hash/fnv"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRemotefile() *schema.Resource {
	return &schema.Resource{
		Description: "Sample data source in the Terraform provider scaffolding.",

		ReadContext: dataSourceRemotefileRead,

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
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("REMOTEFILE_PORT", 22),
				Description: "The ssh port to the target host. May alternatively be set via the `REMOTEFILE_PORT` environment variable.",
			},
			"path": {
				Description: "Path to file on remote host.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceRemotefileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := resourceRemotefileRead(ctx, d, meta)

	if !err.HasError() {
		h := fnv.New32a()
		h.Write([]byte(fmt.Sprintf("%s @ %s", d.Get("content"), d.Get("path"))))
		d.SetId(strconv.Itoa(int(h.Sum32())))
	}

	return err
}
