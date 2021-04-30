package provider

import (
	"bytes"
	"context"
	"strings"

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
						"private_key": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "The private key used to login to the target host.",
						},
					},
				},
			},
			"path": {
				Description: "Path to file on remote host.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"content": {
				Description: "Content of file.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"permissions": {
				Description: "Permissions of file.",
				Type:        schema.TypeString,
				Default:     "0644",
				Optional:    true,
			},
		},
	}
}

func resourceRemotefileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)

	d.SetId(d.Get("path").(string))

	resourceClient, err := client.fromResourceData(d)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	scpClient, err := resourceClient.getSCPClient()
	if err != nil {
		return diag.Errorf(err.Error())
	}
	defer scpClient.Close()

	err = scpClient.CopyFile(strings.NewReader(d.Get("content").(string)), d.Get("path").(string), d.Get("permissions").(string))

	if err != nil {
		return diag.Errorf("error while copying file: %s", err.Error())
	}

	return diag.Diagnostics{}
}

func resourceRemotefileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)

	resourceClient, err := client.fromResourceData(d)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	sftpClient, err := resourceClient.getSFTPClient()
	if err != nil {
		return diag.Errorf(err.Error())
	}
	defer sftpClient.Close()

	file, err := sftpClient.Open(d.Get("path").(string))
	if err != nil {
		return diag.Errorf("error while opening remote file: %s", err.Error())
	}
	defer file.Close()

	content := bytes.Buffer{}
	_, err = file.WriteTo(&content)

	if err != nil {
		return diag.Errorf("error while reading remote file: %s", err.Error())
	}

	d.Set("content", string(content.String()))

	return diag.Diagnostics{}
}

func resourceRemotefileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceRemotefileCreate(ctx, d, meta)
}

func resourceRemotefileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)

	resourceClient, err := client.fromResourceData(d)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	sftpClient, err := resourceClient.getSFTPClient()
	if err != nil {
		return diag.Errorf(err.Error())
	}
	defer sftpClient.Close()

	err = sftpClient.Remove(d.Get("path").(string))
	if err != nil {
		return diag.Errorf("error while removing remote file: %s", err.Error())
	}

	return diag.Diagnostics{}
}
