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
  name        = "Telephony Resource Test"
  description = "Testing Topology Telephony Resource Management"
  notes       = "Created via Terraform Test"
  datacenter  = "LON"
}

resource "dcloud_telephony" "test_telephony" {
  inventory_telephony_id = "1"
  topology_uid         = dcloud_topology.test_topology.id
}