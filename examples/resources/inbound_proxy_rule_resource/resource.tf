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

resource "dcloud_topology" "test_topology" {
  name        = "Inbound Proxy Resource Test"
  description = "Testing Topology Inbound Proxy Resource Management"
  notes       = "Created via Terraform Test"
  datacenter  = "LON"
}

resource "dcloud_network" "unrouted_network" {
  name                 = "An unrouted network"
  description          = "Demonstrating a network not routed through VPOD Gateway"
  inventory_network_id = "L2-VLAN-16"
  topology_uid         = dcloud_topology.test_topology.id
}

resource "dcloud_vm" "vm1" {
  inventory_vm_id = "7668085"
  topology_uid    = dcloud_topology.test_topology.id
  name            = "Ubuntu Desktop 1"
  description     = "A standard Ubuntu Desktop VM"
  cpu_qty         = 1
  memory_mb       = 1024

  network_interfaces {
    network_uid    = dcloud_network.unrouted_network.id
    name           = "Network adapter 1"
    mac_address    = "00:50:56:00:03:AA"
    type           = "VIRTUAL_E1000"
    ip_address     = "127.0.0.2"
    ssh_enabled    = true
    rdp_enabled    = true
    rdp_auto_login = true
  }

  advanced_settings {
    all_disks_non_persistent = false
    bios_uuid                = "42 3a 5f 9d f1 a8 7c 0e-7d c2 44 27 2e d6 67 aa"
    name_in_hypervisor       = "ubuntu"
    not_started              = false
  }

  remote_access {
    vm_console_enabled = true
    display_credentials {
      username = "displayuser"
      password = "displaypassword"
    }
  }

}

resource "dcloud_inbound_proxy_rule" "inbound_proxy_rule" {
  topology_uid = dcloud_topology.test_topology.id
  nic_uid      = dcloud_vm.vm1.network_interfaces[0].uid
  tcp_port     = 443
  url_path     = "/testing/url/"
  hyperlink    = "Click Me"
  ssl          = true
}