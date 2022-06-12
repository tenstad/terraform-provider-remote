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
	content := d.Get("content").(string)
	path := d.Get("path").(string)
	permissions := d.Get("permissions").(string)
	group := d.Get("group").(string)
	owner := d.Get("owner").(string)

	client, err := meta.(*apiClient).getRemoteClient(conn)
	if err != nil {
		return diag.Errorf("error while opening remote client: %s", err.Error())
	}

	sudo, ok := conn.GetOk("conn.0.sudo")
	if ok && sudo.(bool) {
		err := client.WriteFileSudo(content, path)
		if err != nil {
			return diag.Errorf("error while creating remote file with sudo: %s", err.Error())
		}
		err = client.ChmodFile(path, permissions, true)
		if err != nil {
			return diag.Errorf("error while changing permissions of remote file with sudo: %s", err.Error())
		}
		if group != "" {
			err = client.ChgrpFile(path, group, true)
			if err != nil {
				return diag.Errorf("error while changing group of remote file with sudo: %s", err.Error())
			}
		}
		if owner != "" {
			err = client.ChownFile(path, owner, true)
			if err != nil {
				return diag.Errorf("error while changing owner of remote file with sudo: %s", err.Error())
			}
		}
	} else {
		err := client.WriteFile(content, path, permissions)
		if err != nil {
			return diag.Errorf("error while creating remote file: %s", err.Error())
		}
		err = client.ChmodFile(path, permissions, false)
		if err != nil {
			return diag.Errorf("error while changing permissions of remote file: %s", err.Error())
		}
		if group != "" {
			err = client.ChgrpFile(path, group, false)
			if err != nil {
				return diag.Errorf("error while changing group of remote file: %s", err.Error())
			}
		}
		if owner != "" {
			err = client.ChownFile(path, owner, false)
			if err != nil {
				return diag.Errorf("error while changing owner of remote file: %s", err.Error())
			}
		}
	}

	err = meta.(*apiClient).closeRemoteClient(conn)
	if err != nil {
		return diag.Errorf("error while closing remote client: %s", err.Error())
	}

	return diag.Diagnostics{}
}

func resourceRemoteFileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := meta.(*apiClient).getConnWithDefault(d)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	setResourceID(d, conn)
	path := d.Get("path").(string)
	group := d.Get("group").(string)
	owner := d.Get("owner").(string)

	client, err := meta.(*apiClient).getRemoteClient(conn)
	if err != nil {
		return diag.Errorf("error while opening remote client: %s", err.Error())
	}

	sudo, ok := conn.GetOk("conn.0.sudo")
	if ok && sudo.(bool) {
		exists, err := client.FileExistsSudo(path)
		if err != nil {
			return diag.Errorf("error while checking if remote file exists with sudo: %s", err.Error())
		}
		if exists {
			content, err := client.ReadFileSudo(path)
			if err != nil {
				return diag.Errorf("error while reading remote file with sudo: %s", err.Error())
			}
			d.Set("content", content)
			permissions, err := client.ReadFilePermissions(path, true)
			if err != nil {
				return diag.Errorf("error while reading remote file permissions with sudo: %s", err.Error())
			}
			d.Set("permissions", permissions)
			if owner != "" {
				owner, err := client.ReadFileOwner(path, true)
				if err != nil {
					return diag.Errorf("error while reading remote file owner with sudo: %s", err.Error())
				}
				d.Set("owner", owner)
			}
			if group != "" {
				group, err := client.ReadFileGroup(path, true)
				if err != nil {
					return diag.Errorf("error while reading remote file group with sudo: %s", err.Error())
				}
				d.Set("group", group)
			}
		} else {
			return diag.Errorf("cannot read file, it does not exist.")
		}
	} else {
		content, err := client.ReadFile(path)
		if err != nil {
			return diag.Errorf("error while reading remote file: %s", err.Error())
		}
		d.Set("content", content)
		permissions, err := client.ReadFilePermissions(path, false)
		if err != nil {
			return diag.Errorf("error while reading remote file permissions: %s", err.Error())
		}
		d.Set("permissions", permissions)
		if owner != "" {
			owner, err := client.ReadFileOwner(path, false)
			if err != nil {
				return diag.Errorf("error while reading remote file owner: %s", err.Error())
			}
			d.Set("owner", owner)
		}
		if group != "" {
			group, err := client.ReadFileGroup(path, false)
			if err != nil {
				return diag.Errorf("error while reading remote file group: %s", err.Error())
			}
			d.Set("group", group)
		}
	}

	err = meta.(*apiClient).closeRemoteClient(conn)
	if err != nil {
		return diag.Errorf("error while closing remote client: %s", err.Error())
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
		return diag.Errorf("error while opening remote client: %s", err.Error())
	}

	path := d.Get("path").(string)

	sudo, ok := conn.GetOk("conn.0.sudo")
	if ok && sudo.(bool) {
		exists, err := client.FileExistsSudo(path)
		if err != nil {
			return diag.Errorf("error while checking if remote file exists with sudo: %s", err.Error())
		}
		if exists {
			err := client.DeleteFileSudo(path)
			if err != nil {
				return diag.Errorf("error while removing remote file with sudo: %s", err.Error())
			}
		} else {
			return diag.Errorf("cannot delete file, it does not exist.")
		}
	} else {
		err := client.DeleteFile(path)
		if err != nil {
			return diag.Errorf("error while removing remote file: %s", err.Error())
		}
	}

	err = meta.(*apiClient).closeRemoteClient(conn)
	if err != nil {
		return diag.Errorf("error while closing remote client: %s", err.Error())
	}

	return diag.Diagnostics{}
}
