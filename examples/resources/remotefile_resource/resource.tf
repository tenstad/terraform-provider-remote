resource "remotefile" "foo" {
  conn {
    host        = "foo.com"
    port        = 22
    username    = "foo"
    sudo        = true
    private_key = "<ssh private key>"
  }
  path = "/tmp/foo.txt"
  content = "bar"
  permissions = "0777"
}
