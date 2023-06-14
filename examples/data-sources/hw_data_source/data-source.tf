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
  name        = "Test Topology For Testing Hardware"
  description = "Will be used to load hardware"
  notes       = ""
  datacenter  = "LON"
}

resource "dcloudtb_hw" "IE4000" {
  topology_uid               = dcloudtb_topology.test_topology.id
  inventory_hw_id            = "76"
  name                       = "IE 4000 Device"
  hardware_console_enabled   = false
}

data "dcloudtb_hws" "test_topology_hws" {
  topology_uid = dcloudtb_topology.test_topology.id
}

output "hws" {
  value = data.dcloudtb_hws.test_topology_hws
}