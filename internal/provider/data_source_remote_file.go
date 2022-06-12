package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRemoteFile() *schema.Resource {
	return &schema.Resource{
		Description: "File on remote host.",

		ReadContext: dataSourceRemoteFileRead,

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
				Required:    true,
			},
			"content": {
				Description: "Content of file.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"permissions": {
				Description: "Permissions of file (in octal form).",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"group": {
				Description: "Group (GID) of file.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"owner": {
				Description: "Owner (UID) of file.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceRemoteFileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := meta.(*apiClient).getConnWithDefault(d)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	setResourceID(d, conn)
	path := d.Get("path").(string)

	client, err := meta.(*apiClient).getRemoteClient(conn)
	if err != nil {
		return diag.Errorf("unable to open remote client: %s", err.Error())
	}

	sudo, ok := conn.GetOk("conn.0.sudo")
	if ok && sudo.(bool) {
		exists, err := client.FileExistsSudo(path)
		if err != nil {
			return diag.Errorf("unable to check if remote file exists with sudo: %s", err.Error())
		}
		if exists {
			content, err := client.ReadFileSudo(path)
			if err != nil {
				return diag.Errorf("unable to read remote file with sudo: %s", err.Error())
			}
			d.Set("content", content)
			permissions, err := client.ReadFilePermissions(path, true)
			if err != nil {
				return diag.Errorf("unable to read remote file permissions with sudo: %s", err.Error())
			}
			d.Set("permissions", permissions)
			owner, err := client.ReadFileOwner(path, true)
			if err != nil {
				return diag.Errorf("unable to read remote file owner with sudo: %s", err.Error())
			}
			d.Set("owner", owner)
			group, err := client.ReadFileGroup(path, true)
			if err != nil {
				return diag.Errorf("unable to read remote file group with sudo: %s", err.Error())
			}
			d.Set("group", group)
		} else {
			return diag.Errorf("cannot read remote file, it does not exist.")
		}
	} else {
		content, err := client.ReadFile(path)
		if err != nil {
			return diag.Errorf("unable to read remote file: %s", err.Error())
		}
		d.Set("content", content)
		permissions, err := client.ReadFilePermissions(path, false)
		if err != nil {
			return diag.Errorf("unable to read remote file permissions: %s", err.Error())
		}
		d.Set("permissions", permissions)
		owner, err := client.ReadFileOwner(path, false)
		if err != nil {
			return diag.Errorf("unable to read remote file owner: %s", err.Error())
		}
		d.Set("owner", owner)
		group, err := client.ReadFileGroup(path, false)
		if err != nil {
			return diag.Errorf("unable to read remote file group: %s", err.Error())
		}
		d.Set("group", group)
	}

	err = meta.(*apiClient).closeRemoteClient(conn)
	if err != nil {
		return diag.Errorf("unable to close remote client: %s", err.Error())
	}

	return diag.Diagnostics{}
}
