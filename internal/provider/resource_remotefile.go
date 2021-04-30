package provider

import (
	"bytes"
	"context"
	"fmt"
	"hash/fnv"
	"strconv"
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
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("REMOTEFILE_USERNAME", nil),
				Description: "The username on the target host. May alternatively be set via the `REMOTEFILE_USERNAME` environment variable.",
			},
			"private_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("REMOTEFILE_PRIVATE_KEY", nil),
				Description: "The private key used to login to target host. May alternatively be set via the `REMOTEFILE_PRIVATE_KEY` environment variable.",
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

	h := fnv.New32a()
	h.Write([]byte(fmt.Sprintf("%s @ %s", d.Get("content"), d.Get("path"))))
	d.SetId(strconv.Itoa(int(h.Sum32())))

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
