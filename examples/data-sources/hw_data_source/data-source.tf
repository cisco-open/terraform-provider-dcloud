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

data "dcloudtb_hws" "test_topology_hws" {
  topology_uid = "6g9d52l7oq830mhqutsx44qhg"
}

output "hws" {
  value = data.dcloudtb_hws.test_topology_hws
}