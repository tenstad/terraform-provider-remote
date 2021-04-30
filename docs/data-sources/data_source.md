---
page_title: "remotefile_data_source Data Source - terraform-provider-remotefile"
subcategory: ""
description: |-
  Remote file datasource.
---

# Data Source `remotefile_data_source`

  Remote file datasource.

## Example Usage

```terraform
data "remotefile_data_source" "bar" {
	path = "/tmp/bar.txt"
}
```

## Schema

### Required

- **path** (String, Required) Path to file on remote host.

### Optional

- *none*
