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
  name        = "License Resource Test"
  description = "Testing Topology License Resource Management"
  notes       = "Created via Terraform Test"
  datacenter  = "LON"
}

resource "dcloud_license" "test_license" {
  quantity             = 3
  inventory_license_id = "340"
  topology_uid         = dcloud_topology.test_topology.id
}