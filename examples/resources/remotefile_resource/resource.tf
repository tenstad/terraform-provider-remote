resource "remotefile" "foo" {
  path = "/tmp/foo.txt"
  content = "bar"
  permissions = "0777"
}
