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
  name        = "HW Start Order Resource Test"
  description = "Testing Topology HW Start Order Resource Management"
  notes       = "Created via Terraform Test"
  datacenter  = "LON"
}

resource "dcloud_hw" "hw1" {
  topology_uid    = dcloud_topology.test_topology.id
  inventory_hw_id = "76"
  name            = "IE 4000 Device"
}

resource "dcloud_hw" "hw2" {
  topology_uid    = dcloud_topology.test_topology.id
  inventory_hw_id = "14"
  name            = "UCS Hardware Pod"
}

resource "dcloud_hw_start_order" "hw_start_order" {
  ordered      = true
  topology_uid = dcloud_topology.test_topology.id

  start_positions {
    position      = 1
    delay_seconds = 10
    hw_uid        = dcloud_hw.hw1.id
  }

  start_positions {
    position      = 2
    delay_seconds = 20
    hw_uid        = dcloud_hw.hw2.id
  }
}