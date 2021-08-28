data "remote_file" "hosts" {
  conn {
    host     = "10.0.0.12"
    username = "john"
    password = "password"
    sudo     = true
  }
  path = "/etc/hosts"
}
