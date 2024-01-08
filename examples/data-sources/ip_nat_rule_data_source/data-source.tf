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
  name        = "Test Topology For Testing IP NAT Rule"
  description = "Will be used to load IP NAT rules "
  notes       = ""
  datacenter  = "LON"
}

resource "dcloud_ip_nat_rule" "ip_nat"{
  topology_uid = dcloud_topology.test_topology.id
  target_ip_address = "192.168.1.1"
  target_name = "Sample Device"
  east_west = true
}

data "dcloud_ip_nat_rules" "test_topology_ip_nat_rules"{
  topology_uid = dcloud_topology.test_topology.id
  depends_on = [dcloud_ip_nat_rule.ip_nat]
}

output "ip_nat_rules" {
  value = data.dcloud_ip_nat_rules.test_topology_ip_nat_rules
}