provider "remote" {
  max_sessions = 2

  conn {
    host     = "localhost"
    port     = 8022
    user     = "root"
    password = "password"
    sudo     = false
  }
}

resource "remote_file" "foo" {
  path        = "/tmp/.bashrc"
  content     = "alias ll='ls -alF'"
  permissions = "0644"
  owner       = "1001"
  group       = "1001"
}
