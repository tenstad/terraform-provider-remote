---
page_title: "remotefile Provider"
subcategory: ""
description: |-
  Manage files on remote hosts
---

# remotefile Provider

## Example Usage

```terraform
provider "remotefile" {
    max_sessions = 2
}
```

## Schema

### Required

- *none*

### Optional

- **max_sessions** (Number, Optional) Maximum number of open sessions in each host connection.
