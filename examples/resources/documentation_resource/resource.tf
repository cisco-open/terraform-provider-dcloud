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
  name        = "Documentation Resource Test"
  description = "Testing Topology Documentation Resource Management"
  notes       = "Created via Terraform Test"
  datacenter  = "LON"
}

resource "dcloud_documentation" "test_documentation" {
  topology_uid = dcloud_topology.test_topology.id
  doc_url      = "https://johndoe.com"
}