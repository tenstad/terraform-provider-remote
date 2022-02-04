provider "remote" {
  max_sessions = 2

  conn {
    host     = "localhost"
    port     = 8022
    user     = "root"
    password = "password"
  }
}

resource "remote_file" "foo" {
  path        = "/tmp/.bashrc"
  content     = "alias ll='ls -alF'"
  permissions = "0644"
}
