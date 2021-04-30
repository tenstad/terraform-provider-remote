---
page_title: "remotefile Provider"
subcategory: ""
description: |-
  Mangae files on remote hosts
---

# remotefile Provider

## Example Usage

```terraform
provider "scaffolding" {
  username = "foo"
  private_key_path = "~/.ssh/id_rsa"
  host = "foo.com"
  port = "22"
}
```

## Schema
