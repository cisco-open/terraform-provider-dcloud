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
  name        = "Test Topology For Testing External DNS rules"
  description = "Will be used to load External DNS Rules"
  notes       = ""
  datacenter  = "LON"
}

resource "dcloud_ip_nat_rule" "test_topology_ip_nat_rule" {
  topology_uid      = dcloud_topology.test_topology.id
  target_ip_address = "192.168.1.1"
  target_name       = "Sample Device"
  east_west         = false
}

resource "dcloud_external_dns" "test_topology_external_dns" {
  topology_uid = dcloud_topology.test_topology.id
  nat_rule_id  = dcloud_ip_nat_rule.test_topology_ip_nat_rule.id
  hostname     = "localhost"
  srv_records {
    service  = "_test"
    protocol = "TCP"
    port     = 8081
  }
}

data "dcloud_external_dns" "external_dns_test" {
  depends_on   = [dcloud_external_dns.test_topology_external_dns]
  topology_uid = dcloud_topology.test_topology.id
}

output "external_dns" {
  value = data.dcloud_external_dns.external_dns_test
}