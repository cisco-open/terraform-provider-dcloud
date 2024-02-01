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

data "dcloud_inventory_srv_protocols" "test_srv_protocols" {
}

output "srv_protocols" {
  value = data.dcloud_inventory_srv_protocols.test_srv_protocols
}