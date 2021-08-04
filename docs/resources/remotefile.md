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
resource "remotefile" "bashrc" {
  conn {
    host        = "google.com"
    port        = 22
    username    = "john"
    private_key = "<ssh private key>"
  }
  path = "/home/john/.bashrc"
  content = "alias ll='ls -alF'"
  permissions = "0644"
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
- **content** (String, Required) Content of file on remote host.

### Optional

- **permissions** (String, Optional) The file permissions.
