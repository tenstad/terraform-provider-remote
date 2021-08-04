package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRemotefile() *schema.Resource {
	return &schema.Resource{
		Description: "File on remote host.",

		ReadContext: dataSourceRemotefileRead,

		Schema: map[string]*schema.Schema{
			"conn": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Connection to host where files are located.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The target host.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     22,
							Description: "The ssh port to the target host.",
						},
						"username": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The username on the target host.",
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
							Description: "The pasword for the user on the target host.",
						},
						"private_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "The private key used to login to the target host.",
						},
						"private_key_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The local path to the private key used to login to the target host",
						},
						"private_key_env_var": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the local environment variable containing the private key used to login to the target host.",
						},
					},
				},
			},
			"path": {
				Description: "Path to file on remote host.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},
			"content": {
				Description: "Content of file.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceRemotefileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceRemotefileRead(ctx, d, meta)
}
