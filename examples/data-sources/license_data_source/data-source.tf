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
  name        = "Test Topology For Testing Licenses"
  description = "Will be used to load licenses"
  notes       = ""
  datacenter  = "LON"
}

resource "dcloud_license" "test_license" {
  quantity             = 3
  inventory_license_id = "340"
  topology_uid         = dcloud_topology.test_topology.id
}

data "dcloud_licenses" "test_topology_licenses" {
  topology_uid = dcloud_topology.test_topology.id
}

output "licenses" {
  value = data.dcloud_licenses.test_topology_licenses
}