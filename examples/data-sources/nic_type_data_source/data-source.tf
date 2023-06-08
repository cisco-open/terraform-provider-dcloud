terraform {
  required_providers {
    dcloudtb = {
      version = "0.1"
      source  = "cisco.com/dcloud/dcloudtb"
    }
  }
}

provider "dcloudtb" {
  tb_url = "https://tbv3-production.ciscodcloud.com/api"
}

data "dcloudtb_nic_types" "nic_types" {}

output "nic_types" {
  value = data.dcloudtb_nic_types.nic_types
}