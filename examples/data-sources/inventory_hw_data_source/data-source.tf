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
  name        = "Test Topology For Testing Inventory HWs"
  description = "Will be used to load inventory hws"
  notes       = ""
  datacenter  = "LON"
}

data "dcloud_inventory_hws" "test_topology_inventory_hws" {
  topology_uid = dcloud_topology.test_topology.id
}

output "hws" {
  value = data.dcloud_inventory_hws.test_topology_inventory_hws
}