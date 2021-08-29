data "remote_file" "hosts" {
  conn {
    host     = "10.0.0.17"
    user     = "john"
    password = "password"
    sudo     = true
  }

  path = "/etc/hosts"
}

data "remote_file" "server1_hosts" {
  provider = remote.server1

  path = "/etc/hosts"
}

data "remote_file" "server2_hosts" {
  provider = remote.server2

  path = "/etc/hosts"
}
