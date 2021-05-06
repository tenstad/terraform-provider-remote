---
page_title: "remotefile Data Source - terraform-provider-remotefile"
subcategory: ""
description: |-
  Remote file datasource.
---

# Data Source `remotefile`

  Remote file datasource.

## Example Usage

```terraform
data "remotefile" "bar" {
  conn {
    host     = "foo.com"
    port     = 22
    username = "foo"
    password = "<password>"
  }
  path = "/tmp/bar.txt"
}
```

## Schema

### Required

- **conn** (Object, Required) Connection to remote host.
- **path** (String, Required) Path to file on remote host.

### Optional

- *none*
