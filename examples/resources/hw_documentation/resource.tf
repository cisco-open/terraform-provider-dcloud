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
  name        = "My Terraformed Topology"
  description = "A topology created from terraform"
  notes       = "Programmatic clients rule!"
  datacenter  = "LON"
}

resource "dcloud_documentation" "documentation" {
  topology_uid = dcloud_topology.test_topology.id
  doc_url = "https://johndoe.com"
}