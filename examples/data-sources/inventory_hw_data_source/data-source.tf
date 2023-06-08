terraform {
  required_providers {
    dcloudtb = {
      version = "0.1"
      source  = "cisco.com/dcloud/dcloudtb"
    }
  }
}

provider "dcloudtb" {
  tb_url = "https://tbv3-production.ciscodcloud.com/api"
}

resource "dcloudtb_topology" "test_topology" {
  name        = "Test Topology For Testing Inventory HWs"
  description = "Will be used to load inventory hws"
  notes       = ""
  datacenter  = "LON"
}

data "dcloudtb_inventory_hws" "test_topology_inventory_hws" {
  topology_uid = dcloudtb_topology.test_topology.id
}

output "hws" {
  value = data.dcloudtb_inventory_hws.test_topology_inventory_hws
}