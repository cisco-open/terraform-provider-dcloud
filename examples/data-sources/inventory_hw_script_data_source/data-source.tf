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
  name        = "Test Topology For Testing Inventory HW Scripts"
  description = "Will be used to load inventory hw scripts"
  notes       = ""
  datacenter  = "LON"
}

data "dcloud_inventory_hw_scripts" "topology1_inventory_hw_scripts" {
  topology_uid = dcloud_topology.test_topology.id
}

output "scripts" {
  value = data.dcloud_inventory_hw_scripts.topology1_inventory_hw_scripts
}