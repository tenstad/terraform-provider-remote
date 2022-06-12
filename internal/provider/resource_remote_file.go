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
				ForceNew:    true,
				Required:    true,
			},
			"permissions": {
				Description: "Permissions of file (in octal form).",
				Type:        schema.TypeString,
				ForceNew:    false,
				Default:     "0644",
				Optional:    true,
			},
			"group": {
				Description: "Group (GID) of file.",
				Type:        schema.TypeString,
				ForceNew:    false,
				Optional:    true,
			},
			"owner": {
				Description: "Owner (UID) of file.",
				Type:        schema.TypeString,
				ForceNew:    false,
				Optional:    true,
			},
		},
	}
}

func resourceRemoteFileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := meta.(*apiClient).getConnWithDefault(d)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	setResourceID(d, conn)

	client, err := meta.(*apiClient).getRemoteClient(conn)
	if err != nil {
		return diag.Errorf("unable to open remote client: %s", err.Error())
	}

	conn_sudo, ok := conn.GetOk("conn.0.sudo")
	sudo := ok && conn_sudo.(bool)
	content := d.Get("content").(string)
	path := d.Get("path").(string)
	permissions := d.Get("permissions").(string)
	group := d.Get("group").(string)
	owner := d.Get("owner").(string)

	err = client.WriteFile(content, path, permissions, sudo)
	if err != nil {
		return diag.Errorf("unable to create remote file: %s", err.Error())
	}

	err = client.ChmodFile(path, permissions, sudo)
	if err != nil {
		return diag.Errorf("unable to change permissions of remote file: %s", err.Error())
	}

	if group != "" {
		err = client.ChgrpFile(path, group, sudo)
		if err != nil {
			return diag.Errorf("unable to change group of remote file: %s", err.Error())
		}
	}

	if owner != "" {
		err = client.ChownFile(path, owner, sudo)
		if err != nil {
			return diag.Errorf("unable to change owner of remote file: %s", err.Error())
		}
	}

	err = meta.(*apiClient).closeRemoteClient(conn)
	if err != nil {
		return diag.Errorf("unable to close remote client: %s", err.Error())
	}

	return diag.Diagnostics{}
}

func resourceRemoteFileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := meta.(*apiClient).getConnWithDefault(d)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	setResourceID(d, conn)

	client, err := meta.(*apiClient).getRemoteClient(conn)
	if err != nil {
		return diag.Errorf("unable to open remote client: %s", err.Error())
	}

	conn_sudo, ok := conn.GetOk("conn.0.sudo")
	sudo := ok && conn_sudo.(bool)
	path := d.Get("path").(string)
	group := d.Get("group").(string)
	owner := d.Get("owner").(string)

	exists, err := client.FileExists(path, sudo)
	if err != nil {
		return diag.Errorf("unable to check if remote file exists: %s", err.Error())
	}
	if exists {
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

		if owner != "" {
			owner, err := client.ReadFileOwner(path, sudo)
			if err != nil {
				return diag.Errorf("unable to read remote file owner: %s", err.Error())
			}
			d.Set("owner", owner)
		}

		if group != "" {
			group, err := client.ReadFileGroup(path, sudo)
			if err != nil {
				return diag.Errorf("unable to read remote file group: %s", err.Error())
			}
			d.Set("group", group)
		}
	} else {
		d.SetId("")
	}

	err = meta.(*apiClient).closeRemoteClient(conn)
	if err != nil {
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
		return diag.Errorf(err.Error())
	}

	client, err := meta.(*apiClient).getRemoteClient(conn)
	if err != nil {
		return diag.Errorf("unable to open remote client: %s", err.Error())
	}

	conn_sudo, ok := conn.GetOk("conn.0.sudo")
	sudo := ok && conn_sudo.(bool)
	path := d.Get("path").(string)

	exists, err := client.FileExists(path, sudo)
	if err != nil {
		return diag.Errorf("unable to check if remote file exists: %s", err.Error())
	}
	if exists {
		err = client.DeleteFile(path, sudo)
		if err != nil {
			return diag.Errorf("unable to delete remote file: %s", err.Error())
		}
	}

	err = meta.(*apiClient).closeRemoteClient(conn)
	if err != nil {
		return diag.Errorf("unable to close remote client: %s", err.Error())
	}

	return diag.Diagnostics{}
}
