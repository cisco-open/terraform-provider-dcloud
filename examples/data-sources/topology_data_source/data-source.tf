terraform {
  required_providers {
    dcloudtb = {
      version = "0.1"
      source  = "cisco.com/dcloud/dcloudtb"
    }
  }
}

provider "dcloudtb" {
  tb_url = "https://tbv3-dev.dev.ciscodcloud.com/api/"
}

data "dcloudtb_topologies" "all" {}

output "all_topologies" {
  value = data.dcloudtb_topologies.all
}