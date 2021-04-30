---
page_title: "remotefile Resource - terraform-provider-remotefile"
subcategory: ""
description: |-
  File on remote host.
---

# Resource `remotefile`

File on remote host.

## Example Usage

```terraform
resource "remotefile" "foo" {
  path = "/tmp/foo.txt"
  content = "bar"
  permissions = "0777"
}

```

## Schema

### Required

- **path** (String, Required) Path to file on remote host.
- **content** (String, Required) Content of file on remote host.

### Optional

- **permissions** (String, Optional) The file permissions.
