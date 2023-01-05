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
  name        = "Test Topology For Testing Inventory HW Template Configs"
  description = "Will be used to load inventory hw template configs"
  notes       = ""
  datacenter  = "LON"
}

data "dcloudtb_inventory_hw_template_configs" "topology1_inventory_hw_template_configs" {
  topology_uid = dcloudtb_topology.test_topology.id
}

output "scripts" {
  value = data.dcloudtb_inventory_hw_template_configs.topology1_inventory_hw_template_configs
}