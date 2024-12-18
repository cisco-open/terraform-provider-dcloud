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

data "dcloud_evc_modes" "evc_modes" {}

output "evc_modes" {
  value = data.dcloud_evc_modes.evc_modes
}