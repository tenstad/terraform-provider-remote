package provider

import (
	"context"

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
				Required:    true,
			},
			"permissions": {
				Description: "Permissions of file (in octal form).",
				Type:        schema.TypeString,
				Default:     "0644",
				Optional:    true,
			},
			"group": {
				Description: "Group ID (GID) of file owner. Mutually exclusive with `group_name`.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"group_name": {
				Description:   "Group name of file owner. Mutually exclusive with `group`.",
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"group"},
			},
			"owner": {
				Description: "User ID (UID) of file owner. Mutually exclusive with `owner_name`.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"owner_name": {
				Description:   "User name of file owner. Mutually exclusive with `owner`.",
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"owner"},
			},
		},
	}
}

func resourceRemoteFileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := meta.(*apiClient).getConnWithDefault(d)
	if err != nil {
		return diag.Diagnostics{{Severity: diag.Error, Summary: err.Error()}}
	}

	if err := setResourceID(d, conn); err != nil {
		return diag.FromErr(err)
	}

	client, err := meta.(*apiClient).getRemoteClient(ctx, conn)
	if err != nil {
		return diag.Errorf("unable to open remote client: %s", err.Error())
	}

	sudo, _, err := GetOk[bool](conn, "conn.0.sudo")
	if err != nil {
		return diag.Diagnostics{{Severity: diag.Error, Summary: err.Error()}}
	}

	content, err := Get[string](d, "content")
	if err != nil {
		return diag.Diagnostics{{Severity: diag.Error, Summary: err.Error()}}
	}

	path, err := Get[string](d, "path")
	if err != nil {
		return diag.Diagnostics{{Severity: diag.Error, Summary: err.Error()}}
	}

	permissions, err := Get[string](d, "permissions")
	if err != nil {
		return diag.Diagnostics{{Severity: diag.Error, Summary: err.Error()}}
	}

	var group string
	if g, ok, err := GetOk[string](d, "group"); ok {
		if err != nil {
			return diag.Diagnostics{{Severity: diag.Error, Summary: err.Error()}}
		}
		group = g
	} else if g, ok, err := GetOk[string](d, "group_name"); ok {
		if err != nil {
			return diag.Diagnostics{{Severity: diag.Error, Summary: err.Error()}}
		}
		group = g
	}

	var owner string
	if o, ok, err := GetOk[string](d, "owner"); ok {
		if err != nil {
			return diag.Diagnostics{{Severity: diag.Error, Summary: err.Error()}}
		}
		owner = o
	} else if o, ok, err := GetOk[string](d, "owner_name"); ok {
		if err != nil {
			return diag.Diagnostics{{Severity: diag.Error, Summary: err.Error()}}
		}
		owner = o
	}

	if d.HasChange("content") {
		if err := client.WriteFile(ctx, content, path, permissions, sudo); err != nil {
			return diag.Errorf("unable to create remote file: %s", err.Error())
		}
	}

	if err := client.ChmodFile(path, permissions, sudo); err != nil {
		return diag.Errorf("unable to change permissions of remote file: %s", err.Error())
	}

	if group != "" {
		if err := client.ChgrpFile(path, group, sudo); err != nil {
			return diag.Errorf("unable to change group of remote file: %s", err.Error())
		}
	}

	if owner != "" {
		if err := client.ChownFile(path, owner, sudo); err != nil {
			return diag.Errorf("unable to change owner of remote file: %s", err.Error())
		}
	}

	if err := meta.(*apiClient).closeRemoteClient(conn); err != nil {
		return diag.Errorf("unable to close remote client: %s", err.Error())
	}

	return diag.Diagnostics{}
}

func resourceRemoteFileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := meta.(*apiClient).getConnWithDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setResourceID(d, conn); err != nil {
		return diag.FromErr(err)
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

	_, groupOk, err := GetOk[string](d, "group")
	if err != nil {
		return diag.Diagnostics{{Severity: diag.Error, Summary: err.Error()}}
	}

	_, groupNameOk, err := GetOk[string](d, "group_name")
	if err != nil {
		return diag.Diagnostics{{Severity: diag.Error, Summary: err.Error()}}
	}

	_, ownerOk, err := GetOk[string](d, "owner")
	if err != nil {
		return diag.Diagnostics{{Severity: diag.Error, Summary: err.Error()}}
	}

	_, ownerNameOk, err := GetOk[string](d, "owner_name")
	if err != nil {
		return diag.Diagnostics{{Severity: diag.Error, Summary: err.Error()}}
	}

	exists, err := client.FileExists(path, sudo)
	if err != nil {
		return diag.Errorf("unable to check if remote file exists: %s", err.Error())
	}
	if exists {
		content, err := client.ReadFile(path, sudo)
		if err != nil {
			return diag.Errorf("unable to read remote file: %s", err.Error())
		}
		if err := d.Set("content", content); err != nil {
			return diag.FromErr(err)
		}

		permissions, err := client.ReadFilePermissions(path, sudo)
		if err != nil {
			return diag.Errorf("unable to read remote file permissions: %s", err.Error())
		}
		if err := d.Set("permissions", permissions); err != nil {
			return diag.FromErr(err)
		}

		if ownerOk {
			owner, err := client.ReadFileOwner(path, sudo)
			if err != nil {
				return diag.Errorf("unable to read remote file owner: %s", err.Error())
			}
			if err := d.Set("owner", owner); err != nil {
				return diag.FromErr(err)
			}
		}
		if ownerNameOk {
			ownerName, err := client.ReadFileOwnerName(path, sudo)
			if err != nil {
				return diag.Errorf("unable to read remote file owner_name: %s", err.Error())
			}
			if err := d.Set("owner_name", ownerName); err != nil {
				return diag.FromErr(err)
			}
		}

		if groupOk {
			group, err := client.ReadFileGroup(path, sudo)
			if err != nil {
				return diag.Errorf("unable to read remote file group: %s", err.Error())
			}
			if err := d.Set("group", group); err != nil {
				return diag.FromErr(err)
			}
		}
		if groupNameOk {
			groupName, err := client.ReadFileGroupName(path, sudo)
			if err != nil {
				return diag.Errorf("unable to read remote file group_name: %s", err.Error())
			}
			if err := d.Set("group_name", groupName); err != nil {
				return diag.FromErr(err)
			}
		}
	} else {
		d.SetId("")
	}

	if err := meta.(*apiClient).closeRemoteClient(conn); err != nil {
		return diag.Errorf("unable to close remote client: %s", err.Error())
	}

	return diag.Diagnostics{}
}

func resourceRemoteFileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceRemoteFileCreate(ctx, d, meta)
}

func resourceRemoteFileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := meta.(*apiClient).getConnWithDefault(d)
	if err != nil {
		return diag.Diagnostics{{Severity: diag.Error, Summary: err.Error()}}
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
	if exists {
		if err := client.DeleteFile(path, sudo); err != nil {
			return diag.Errorf("unable to delete remote file: %s", err.Error())
		}
	}

	if err := meta.(*apiClient).closeRemoteClient(conn); err != nil {
		return diag.Errorf("unable to close remote client: %s", err.Error())
	}

	return diag.Diagnostics{}
}
