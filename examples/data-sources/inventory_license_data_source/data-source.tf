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
  name        = "Test Topology For Testing Inventory Licenses"
  description = "Will be used to load inventory licenses"
  notes       = ""
  datacenter  = "LON"
}

data "dcloud_inventory_licenses" "topology1_inventory_licenses" {
  topology_uid = dcloud_topology.test_topology.id
}

output "licenses" {
  value = data.dcloud_inventory_licenses.topology1_inventory_licenses
}