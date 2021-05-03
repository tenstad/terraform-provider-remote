data "remotefile" "bar" {
  conn {
    host     = "foo.com"
    port     = 22
    username = "foo"
    password = "<password>"
  }
  path = "/tmp/bar.txt"
}
