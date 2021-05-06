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
  conn {
    host        = "foo.com"
    port        = 22
    username    = "foo"
    sudo        = true
    private_key = "<ssh private key>"
  }
  path = "/tmp/foo.txt"
  content = "bar"
  permissions = "0777"
}

```

## Schema

### Required

- **conn** (Object, Required) Connection to remote host.
- **path** (String, Required) Path to file on remote host.
- **content** (String, Required) Content of file on remote host.

### Optional

- **permissions** (String, Optional) The file permissions.
