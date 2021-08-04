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
