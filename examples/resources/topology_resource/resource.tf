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

resource "dcloudtb_topology" "tb1" {
  name = "My First Terraformed Topology"
  description = "The first, and maybe the best"
  notes = "This is a much better note"
  datacenter = "LON"
}

resource "dcloudtb_topology" "tb2" {
  name = "My Second Terraformed Topology"
  description = "The second, it's still pretty ok"
  notes = "This is a much better note"
  datacenter = "RTP"
}

resource "dcloudtb_topology" "tb3" {
  name = "My Third Terraformed Topology"
  description = "The third, sure it's grand"
  notes = "This is a much better note"
  datacenter = "SNG"
}

#output "hugh_topology" {
#  value = dcloudtb_topology.hughtopology
#}