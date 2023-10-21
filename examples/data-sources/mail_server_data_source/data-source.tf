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
  name        = "Test Topology For Testing Mail Servers"
  description = "Will be used to load Mail Servers"
  notes       = ""
  datacenter  = "LON"
}

data "dcloud_mail_servers" "test_topology_mail_server"{
  topology_uid = dcloud_topology.test_topology.id
}

output "mail_servers" {
  value = data.dcloud_mail_servers.test_topology_mail_server
}