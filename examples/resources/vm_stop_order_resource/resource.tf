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
  name        = "VM Stop Order Resource Test"
  description = "Testing Topology VM Stop Order Resource Management"
  notes       = "Created via Terraform Test"
  datacenter  = "LON"
}

resource "dcloud_network" "routed_network" {
  name                 = "A routed network"
  description          = "A Network to attach VMs"
  inventory_network_id = "L3-VLAN-2"
  topology_uid         = dcloud_topology.test_topology.id
}

resource "dcloud_vm" "vm1" {
  inventory_vm_id   = "7048155"
  topology_uid      = dcloud_topology.test_topology.id
  name              = "VM Created from Terraform"
  description       = "It's Alive"
  cpu_qty           = 2
  memory_mb         = 4096
  nested_hypervisor = false
  os_family         = "LINUX"

  advanced_settings {
    all_disks_non_persistent = false
    bios_uuid                = "42 3a 5f 9d f1 a8 7c 0e-7d c2 44 27 2e d6 67 aa"
    name_in_hypervisor       = "cmm"
    not_started              = false
  }

  network_interfaces {
    network_uid = dcloud_network.routed_network.id
    name        = "Network adapter 0"
    mac_address = "00:50:56:00:01:AA"
    type        = "VIRTUAL_E1000"
  }

  remote_access {
    vm_console_enabled = false
  }
}

resource "dcloud_vm" "vm2" {
  inventory_vm_id   = "7182268"
  topology_uid      = dcloud_topology.test_topology.id
  name              = "2nd VM Created from Terraform"
  description       = "It's Alive"
  cpu_qty           = 1
  memory_mb         = 3072
  nested_hypervisor = false
  os_family         = "WINDOWS"

  advanced_settings {
    all_disks_non_persistent = false
    bios_uuid                = "42 3a 5f 9d f1 a8 7c 0e-7d c2 44 27 2e d6 67 bb"
    name_in_hypervisor       = "computer"
    not_started              = false
  }

  network_interfaces {
    network_uid = dcloud_network.routed_network.id
    name        = "Network adapter 0"
    mac_address = "00:50:56:00:02:AA"
    type        = "VIRTUAL_E1000"
  }

  remote_access {
    vm_console_enabled = false
  }
}

resource "dcloud_vm_stop_order" "vm_stop_order" {
  ordered      = true
  topology_uid = dcloud_topology.test_topology.id

  stop_positions {
    position = 1
    vm_uid   = dcloud_vm.vm1.id
  }

  stop_positions {
    position = 2
    vm_uid   = dcloud_vm.vm2.id
  }
}