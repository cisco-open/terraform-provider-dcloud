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
  name        = "Test Topology For Testing Inventory Networks"
  description = "Will be used to load inventory networks"
  notes       = ""
  datacenter  = "LON"
}

data "dcloud_inventory_networks" "topology1_inventory_networks" {
  topology_uid = dcloud_topology.test_topology.id
}

output "networks" {
  value = data.dcloud_inventory_networks.topology1_inventory_networks
}