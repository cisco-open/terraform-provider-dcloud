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

resource "dcloud_topology" "test_topology" {
  name        = "Test Topology For Testing External DNS Assets"
  description = "Will be used to load External DNS assets"
  notes       = ""
  datacenter  = "LON"
}

data "dcloud_inventory_dns_assets" "test_dns_assets" {
  topology_uid = dcloud_topology.test_topology.id
}

output "dns_assets" {
  value = data.dcloud_inventory_dns_assets.test_dns_assets
}