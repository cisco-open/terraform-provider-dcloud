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
  name        = "Test Topology For Testing Telephony Items"
  description = "Will be used to load telephony items"
  notes       = ""
  datacenter  = "LON"
}

resource "dcloud_telephony" "test_telephony" {
  topology_uid = dcloud_topology.test_topology.id
  inventory_telephony_id = "1"
}

data "dcloud_telephony" "test_topology_telephony" {
  topology_uid = dcloud_topology.test_topology.id
}

output "telephony-items" {
  value = data.dcloud_telephony.test_topology_telephony
}