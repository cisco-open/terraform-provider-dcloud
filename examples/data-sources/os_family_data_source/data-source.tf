terraform {
  required_providers {
    dcloudtb = {
      version = "0.1"
      source  = "cisco.com/dcloud/dcloudtb"
    }
  }
}

provider "dcloudtb" {
  tb_url = "https://tbv3-dev.dev.ciscodcloud.com/api"
}

data "dcloudtb_os_families" "os_families" {}

output "os_families" {
  value = data.dcloudtb_os_families.os_families
}