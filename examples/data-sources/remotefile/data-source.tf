data "remotefile" "hosts" {
  conn {
    host     = "google.com"
    username = "john"
    password = "password"
    sudo     = true
  }
  path = "/etc/hosts"
}
