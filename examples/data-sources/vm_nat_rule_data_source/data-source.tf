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
  name        = "Test Topology For Testing VM NAT Rule"
  description = "Will be used to create VM NAT rules "
  notes       = ""
  datacenter  = "LON"
}

resource "dcloud_network" "routed_network" {
  name                 = "A routed network"
  description          = "Demonstrating a network routed through VPOD Gateway"
  inventory_network_id = "L3-VLAN-2"
  topology_uid         = dcloud_topology.test_topology.id
}

resource "dcloud_network" "unrouted_network" {
  name                 = "An unrouted network"
  description          = "Demonstrating a network not routed through VPOD Gateway"
  inventory_network_id = "L2-VLAN-16"
  topology_uid         = dcloud_topology.test_topology.id
}

resource "dcloud_vm" "vm1" {
  inventory_vm_id   = "7668085"
  topology_uid      = dcloud_topology.test_topology.id
  name              = "Ubuntu Desktop 1"
  description       = "A standard Ubuntu Desktop VM"
  cpu_qty           = 8
  memory_mb         = 8192
  nested_hypervisor = false
  os_family         = "LINUX"

  advanced_settings {
    all_disks_non_persistent = false
    bios_uuid                = "42 3a 5f 9d f1 a8 7c 0e-7d c2 44 27 2e d6 67 aa"
    name_in_hypervisor       = "ubuntu"
    not_started              = false
  }

  network_interfaces {
    network_uid = dcloud_network.routed_network.id
    name        = "Network adapter 0"
    mac_address = "00:50:56:00:01:AA"
    type        = "VIRTUAL_E1000"
  }

  network_interfaces{
    network_uid    = dcloud_network.unrouted_network.id
    name           = "Network adapter 1"
    mac_address    = "00:50:56:00:01:AB"
    type           = "VIRTUAL_E1000"
    ip_address     = "127.0.0.2"
    ssh_enabled    = true
    rdp_enabled    = true
    rdp_auto_login = true
  }

  remote_access {
    username           = "user"
    password           = "password"
    vm_console_enabled = true

    display_credentials {
      username = "displayuser"
      password = "displaypassword"
    }
  }

  guest_automation {
    command       = "RUN PROGRAM"
    delay_seconds = 10
  }
}

resource "dcloud_vm_nat_rule" "vm_nat"{
  topology_uid = dcloud_topology.test_topology.id
  nic_uid = dcloud_vm.vm1.network_interfaces[1].uid
  east_west = true
}

data "dcloud_vm_nat_rules" "test_topology_vm_nat_rules"{
  topology_uid = dcloud_topology.test_topology.id
}

output "vm_nat_rules" {
  value = data.dcloud_vm_nat_rules.test_topology_vm_nat_rules
}