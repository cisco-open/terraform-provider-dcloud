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

data "dcloud_nic_types" "nic_types" {}

output "nic_types" {
  value = data.dcloud_nic_types.nic_types
}