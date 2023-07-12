terraform {
  required_providers {
    dcloud = {
      version = "0.1"
      source  = "cisco-open/dcloud"
    }
  }
}

provider "dcloud" {
  tb_url = "https://tbv3-production.ciscodcloud.com/api"
}

resource "dcloud_topology" "tb1" {
  name        = "My First Terraformed Topology"
  description = "The first, and maybe the best"
  notes       = "This is a much better note"
  datacenter  = "LON"
}

resource "dcloud_topology" "tb2" {
  name        = "My Second Terraformed Topology"
  description = "The second, it's still pretty ok"
  notes       = "This is a much better note"
  datacenter  = "RTP"
}

resource "dcloud_topology" "tb3" {
  name        = "My Third Terraformed Topology"
  description = "The third, sure it's grand"
  notes       = "This is a much better note"
  datacenter  = "SNG"
}