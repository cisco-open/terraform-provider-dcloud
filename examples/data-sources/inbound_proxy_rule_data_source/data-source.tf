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
  name        = "Test Topology For Testing Inbound Proxy rules"
  description = "Will be used to create Inbound proxy Rules"
  notes       = ""
  datacenter  = "LON"
}

data "dcloud_inbound_proxy_rules" "test_topology_inbound_proxy" {
  topology_uid = dcloud_topology.test_topology.id
}

output "inbound_proxy_rules" {
  value = data.dcloud_inbound_proxy_rules.test_topology_inbound_proxy
}