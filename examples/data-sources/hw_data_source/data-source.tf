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
  name        = "Test Topology For Testing Hardware"
  description = "Will be used to load hardware"
  notes       = ""
  datacenter  = "LON"
}

resource "dcloud_hw" "IE4000" {
  topology_uid             = dcloud_topology.test_topology.id
  inventory_hw_id          = "76"
  name                     = "IE 4000 Device"
  hardware_console_enabled = false
}

data "dcloud_hws" "test_topology_hws" {
  topology_uid = dcloud_topology.test_topology.id
}

output "hws" {
  value = data.dcloud_hws.test_topology_hws
}