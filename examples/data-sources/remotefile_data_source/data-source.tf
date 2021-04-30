data "remotefile" "bar" {
	conn {
		host        = "foo.com"
		port        = "22"
		username    = "foo"
		private_key = "<ssh private key>"
	}
	path = "/tmp/bar.txt"
}
