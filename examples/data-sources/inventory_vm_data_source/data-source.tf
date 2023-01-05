terraform {
  required_providers {
    dcloudtb = {
      version = "0.1"
      source  = "cisco.com/dcloud/dcloudtb"
    }
  }
}

provider "dcloudtb" {
  tb_url = "https://tbv3-dev.dev.ciscodcloud.com/api"
}

resource "dcloudtb_topology" "test_topology" {
  name        = "Test Topology For Testing Inventory VMs"
  description = "Will be used to load inventory vms"
  notes       = ""
  datacenter  = "LON"
}

data "dcloudtb_inventory_vms" "test_topology_inventory_vms" {
  topology_uid = dcloudtb_topology.test_topology.id
}

output "vms" {
  value = data.dcloudtb_inventory_vms.test_topology_inventory_vms
}