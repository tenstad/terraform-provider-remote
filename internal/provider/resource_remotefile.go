package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRemotefile() *schema.Resource {
	return &schema.Resource{
		Description: "File on remote host.",

		CreateContext: resourceRemotefileCreate,
		ReadContext:   resourceRemotefileRead,
		UpdateContext: resourceRemotefileUpdate,
		DeleteContext: resourceRemotefileDelete,

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
							Description: "Use sudo to gain access to read/write file.",
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
							Description: "The path of the private key used to login to the target host.",
						},
						"private_key_env_var": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The env var of the private key used to login to the target host.",
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
				ForceNew:    true,
				Required:    true,
			},
			"permissions": {
				Description: "Permissions of file.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Default:     "0644",
				Optional:    true,
			},
		},
	}
}

func resourceRemotefileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(fmt.Sprintf("%s:%s", d.Get("conn.0.host").(string), d.Get("path").(string)))

	client, err := meta.(*apiClient).getRemoteClient(d)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	sudo, ok := d.GetOk("conn.0.sudo")
	if ok && sudo.(bool) {
		err := client.WriteFileSudo(d)
		if err != nil {
			return diag.Errorf("error while creating remote file with sudo: %s", err.Error())
		}
		err = client.ChmodFileSudo(d)
		if err != nil {
			return diag.Errorf("error while changing permissions of remote file with sudo: %s", err.Error())
		}
	} else {
		err := client.WriteFile(d)
		if err != nil {
			return diag.Errorf("error while creating remote file: %s", err.Error())
		}
	}

	return diag.Diagnostics{}
}

func resourceRemotefileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(fmt.Sprintf("%s:%s", d.Get("conn.0.host").(string), d.Get("path").(string)))

	client, err := meta.(*apiClient).getRemoteClient(d)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	sudo, ok := d.GetOk("conn.0.sudo")
	if ok && sudo.(bool) {
		exists, err := client.FileExistsSudo(d)
		if err != nil {
			return diag.Errorf("error while checking if remote file exists with sudo: %s", err.Error())
		}
		if exists {
			err := client.ReadFileSudo(d)
			if err != nil {
				return diag.Errorf("error while reading remote file with sudo: %s", err.Error())
			}
		} else {
			return diag.Errorf("cannot read file, it does not exist.")
		}
	} else {
		err := client.ReadFile(d)
		if err != nil {
			return diag.Errorf("error while reading remote file: %s", err.Error())
		}
	}

	return diag.Diagnostics{}
}

func resourceRemotefileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceRemotefileCreate(ctx, d, meta)
}

func resourceRemotefileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*apiClient).getRemoteClient(d)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	sudo, ok := d.GetOk("conn.0.sudo")
	if ok && sudo.(bool) {
		exists, err := client.FileExistsSudo(d)
		if err != nil {
			return diag.Errorf("error while checking if remote file exists with sudo: %s", err.Error())
		}
		if exists {
			err := client.DeleteFileSudo(d)
			if err != nil {
				return diag.Errorf("error while removing remote file with sudo: %s", err.Error())
			}
		} else {
			return diag.Errorf("cannot delete file, it does not exist.")
		}
	} else {
		err := client.DeleteFile(d)
		if err != nil {
			return diag.Errorf("error while removing remote file: %s", err.Error())
		}
	}

	return diag.Diagnostics{}
}
