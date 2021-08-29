resource "remote_file" "bashrc" {
  conn {
    host        = "10.0.0.12"
    port        = 22
    user        = "john"
    private_key = "<ssh private key>"
  }

  path        = "/home/john/.bashrc"
  content     = "alias ll='ls -alF'"
  permissions = "0644"
}

resource "remote_file" "server1_bashrc" {
  provider = remote.server1

  path        = "/home/john/.bashrc"
  content     = "alias ll='ls -alF'"
  permissions = "0644"
}

resource "remote_file" "server2_bashrc" {
  provider = remote.server2

  path        = "/home/john/.bashrc"
  content     = "alias ll='ls -alF'"
  permissions = "0644"
}
