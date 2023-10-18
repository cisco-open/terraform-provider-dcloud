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
  name        = "IP NAT Rule Resource Test"
  description = "Testing IP NAT Rule Resource Management"
  notes       = "Created via Terraform Test"
  datacenter  = "LON"
}

resource "dcloud_ip_nat_rule" "test_ip_nat"{
  topology_uid = dcloud_topology.test_topology.id
  target_ip_address = "192.168.1.1"
  target_name = "Sample Device"
  east_west = true
}