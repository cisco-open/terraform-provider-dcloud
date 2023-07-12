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
  name        = "Test Topology For Testing Inventory VMs"
  description = "Will be used to load inventory vms"
  notes       = ""
  datacenter  = "LON"
}

data "dcloud_inventory_vms" "test_topology_inventory_vms" {
  topology_uid = dcloud_topology.test_topology.id
}

output "vms" {
  value = data.dcloud_inventory_vms.test_topology_inventory_vms
}