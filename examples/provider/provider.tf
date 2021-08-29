# 'conn' must be defined in the resources and data sources when not defined in the provider.
provider "remote" {
  max_sessions = 2
}

# When 'conn' is defined in the provider, it can be overridden in resources and data sources.
# To override it, simply define a 'conn' in the resource or data source.
provider "remote" {
  alias = "server1"

  max_sessions = 2

  conn {
    host     = "10.0.0.12"
    user     = "john"
    password = "password"
    sudo     = true
  }
}

provider "remote" {
  alias = "server2"

  max_sessions = 2

  conn {
    host     = "10.0.0.15"
    user     = "john"
    password = "password"
    sudo     = true
  }
}
