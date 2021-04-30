package provider

import (
	"context"
	"fmt"
	"hash/fnv"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRemotefile() *schema.Resource {
	return &schema.Resource{
		Description: "Sample data source in the Terraform provider scaffolding.",

		ReadContext: dataSourceRemotefileRead,

		Schema: map[string]*schema.Schema{
			"path": {
				Description: "Path to file on remote host.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceRemotefileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := resourceRemotefileRead(ctx, d, meta)

	if !err.HasError() {
		h := fnv.New32a()
		h.Write([]byte(fmt.Sprintf("%s @ %s", d.Get("content"), d.Get("path"))))
		d.SetId(strconv.Itoa(int(h.Sum32())))
	}

	return err
}
