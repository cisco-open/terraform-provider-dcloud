terraform {
  required_providers {
    dcloudtb = {
      version = "0.1"
      source  = "cisco.com/dcloud/dcloudtb"
    }
  }
}

provider "dcloudtb" {
  tb_url = "https://tbv3-dev.dev.ciscodcloud.com/api"
}

resource "dcloudtb_topology" "test_topology" {
  name = "Network Resource Test"
  description = "Testing Topology Network Resource Management"
  notes = "Created via Terraform Test"
  datacenter = "SNG"
}

resource "dcloudtb_network" "routed_network" {
  name = "A routed network"
  description = "Demonstrating a network routed through VPOD Gateway"
  inventory_network_id = "L3-VLAN-2"
  topology_uid = dcloudtb_topology.test_topology.id
}

resource "dcloudtb_network" "unrouted_network" {
  name = "An unrouted network"
  description = "Demonstrating a network, not routed through VPOD Gateway"
  inventory_network_id = "L2-VLAN-16"
  topology_uid = dcloudtb_topology.test_topology.id
}

#data "dcloudtb_networks" "test_topology_networks" {
#  topology_uid = dcloudtb_topology.test_topology.id
#}
#
#output "topology_networks" {
#  value = data.dcloudtb_networks.test_topology_networks
#}