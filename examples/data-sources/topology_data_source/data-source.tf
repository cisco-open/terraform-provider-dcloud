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

data "dcloud_topologies" "all" {}

output "all_topologies" {
  value = data.dcloud_topologies.all
}