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
  debug  = true
}

resource "dcloud_topology" "test_topology" {
  name        = "Scenario Resource Test"
  description = "Testing Topology Scenario Resource Management"
  notes       = "Created via Terraform Test"
  datacenter  = "LON"
}

resource "dcloud_scenario" "scenario" {
  enabled      = true
  question     = "What would you like to do today?"
  topology_uid = dcloud_topology.test_topology.id

  options {
    internal_name = "option1"
    display_name  = "Launch Demo Context1"
  }

  options {
    internal_name = "option2"
    display_name  = "Launch Demo Context2"
  }
}