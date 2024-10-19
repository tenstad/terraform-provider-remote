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
				Description: "Group ID (GID) of file owner.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"group_name": {
				Description: "Group name of file owner.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"owner": {
				Description: "User ID (UID) of file owner.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"owner_name": {
				Description: "User name of file owner.",
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

	if err := setResourceID(d, conn); err != nil {
		return diag.Errorf(err.Error())
	}

	client, err := meta.(*apiClient).getRemoteClient(ctx, conn)
	if err != nil {
		return diag.Errorf("unable to open remote client: %s", err.Error())
	}

	sudo, _, err := GetOk[bool](conn, "conn.0.sudo")
	if err != nil {
		return diag.Diagnostics{{Severity: diag.Error, Summary: err.Error()}}
	}

	path, err := Get[string](d, "path")
	if err != nil {
		return diag.Diagnostics{{Severity: diag.Error, Summary: err.Error()}}
	}

	exists, err := client.FileExists(path, sudo)
	if err != nil {
		return diag.Errorf("unable to check if remote file exists: %s", err.Error())
	}
	if !exists {
		return diag.Errorf("cannot read file, it does not exist")
	}

	content, err := client.ReadFile(path, sudo)
	if err != nil {
		return diag.Errorf("unable to read remote file: %s", err.Error())
	}
	d.Set("content", content)

	permissions, err := client.ReadFilePermissions(path, sudo)
	if err != nil {
		return diag.Errorf("unable to read remote file permissions: %s", err.Error())
	}
	d.Set("permissions", permissions)

	owner, err := client.ReadFileOwner(path, sudo)
	if err != nil {
		return diag.Errorf("unable to read remote file owner: %s", err.Error())
	}
	d.Set("owner", owner)

	owner_name, err := client.ReadFileOwnerName(path, sudo)
	if err != nil {
		return diag.Errorf("unable to read remote file owner_name: %s", err.Error())
	}
	d.Set("owner_name", owner_name)

	group, err := client.ReadFileGroup(path, sudo)
	if err != nil {
		return diag.Errorf("unable to read remote file group: %s", err.Error())
	}
	d.Set("group", group)

	group_name, err := client.ReadFileGroupName(path, sudo)
	if err != nil {
		return diag.Errorf("unable to read remote file group_name: %s", err.Error())
	}
	d.Set("group_name", group_name)

	if err := meta.(*apiClient).closeRemoteClient(conn); err != nil {
		return diag.Errorf("unable to close remote client: %s", err.Error())
	}

	return diag.Diagnostics{}
}
