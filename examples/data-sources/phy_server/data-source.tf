terraform {
  required_providers {
    phy = {
      source  = "registry.terraform.io/sacloud/phy"
    }
  }
}

data "phy_server" "example" {
  filter = "server"
}
