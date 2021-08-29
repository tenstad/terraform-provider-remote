package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRemoteFile() *schema.Resource {
	return &schema.Resource{
		Description: "File on remote host.",

		CreateContext: resourceRemoteFileCreate,
		ReadContext:   resourceRemoteFileRead,
		UpdateContext: resourceRemoteFileUpdate,
		DeleteContext: resourceRemoteFileDelete,

		Schema: map[string]*schema.Schema{
			"conn": {
				Type:        schema.TypeList,
				MinItems:    0,
				MaxItems:    1,
				Optional:    true,
				Description: "Connection to host where files are located.",
				Elem:        connectionSchemaResource,
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

func resourceRemoteFileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(fmt.Sprintf("%s:%s", d.Get("conn.0.host").(string), d.Get("path").(string)))

	client, err := meta.(*apiClient).getRemoteClient(d)
	if err != nil {
		return diag.Errorf("error while opening remote client: %s", err.Error())
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

	err = meta.(*apiClient).closeRemoteClient(d)
	if err != nil {
		return diag.Errorf("error while closing remote client: %s", err.Error())
	}

	return diag.Diagnostics{}
}

func resourceRemoteFileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(fmt.Sprintf("%s:%s", d.Get("conn.0.host").(string), d.Get("path").(string)))

	client, err := meta.(*apiClient).getRemoteClient(d)
	if err != nil {
		return diag.Errorf("error while opening remote client: %s", err.Error())
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

	err = meta.(*apiClient).closeRemoteClient(d)
	if err != nil {
		return diag.Errorf("error while closing remote client: %s", err.Error())
	}

	return diag.Diagnostics{}
}

func resourceRemoteFileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceRemoteFileCreate(ctx, d, meta)
}

func resourceRemoteFileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*apiClient).getRemoteClient(d)
	if err != nil {
		return diag.Errorf("error while opening remote client: %s", err.Error())
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

	err = meta.(*apiClient).closeRemoteClient(d)
	if err != nil {
		return diag.Errorf("error while closing remote client: %s", err.Error())
	}

	return diag.Diagnostics{}
}
