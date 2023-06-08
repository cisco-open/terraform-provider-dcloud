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

data "dcloudtb_vms" "test_topology_vms" {
  topology_uid = "6moevn0brurtqksmystmk8qcb"
}

output "vms" {
  value = data.dcloudtb_vms.test_topology_vms
}