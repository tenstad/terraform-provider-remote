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
data "remotefile" "hosts" {
  conn {
    host     = "google.com"
    username = "john"
    password = "password"
    sudo     = true
  }
  path = "/etc/hosts"
}
```

## Schema

### Required

- **conn** (Object, Required) Connection to remote host.
  - **host** (String, Required) The target host.
  - **port** (Number, Optional) The ssh port to the target host.
  - **username** (String, Required) The username on the target host.
  - **sudo** (Boolean, Optional) Use sudo to gain access to file.
  - **password** (String, Optional) The pasword for the user on the target host.
  - **private_key** (String, Optional) The private key used to login to the target host.
  - **private_key_path** (String, Optional) The local path to the private key used to login to the target host.
  - **private_key_env_var** (String, Optional) The name of the local environment variable containing the private key used to login to the target host.
- **path** (String, Required) Path to file on remote host.

### Optional

- *none*
