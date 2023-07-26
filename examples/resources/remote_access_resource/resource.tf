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
  name        = "Remote Access Resource Test"
  description = "Testing Topology Remote Access Resource Management"
  notes       = "Created via Terraform Test"
  datacenter  = "LON"
}

resource "dcloud_remote_access" "remote_access" {
  any_connect_enabled  = true
  endpoint_kit_enabled = true
  topology_uid         = dcloud_topology.test_topology.id
}