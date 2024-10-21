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
  name        = "Test Topology For Testing Virtual Machines"
  description = "Will be used to load VMs"
  notes       = ""
  datacenter  = "LON"
}

resource "dcloud_network" "routed_network" {
  name                 = "A routed network"
  description          = "Demonstrating a network routed through VPOD Gateway"
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

  network_interfaces {
    network_uid = dcloud_network.routed_network.id
    name        = "Network adapter 0"
    mac_address = "00:50:56:00:01:AA"
    type        = "VIRTUAL_E1000"
  }

  advanced_settings {
    all_disks_non_persistent = false
    bios_uuid                = "42 3a 5f 9d f1 a8 7c 0e-7d c2 44 27 2e d6 67 aa"
    name_in_hypervisor       = "cmm"
    not_started              = false
    evc_mode                 = "SKYLAKE"
  }

  remote_access {
    username           = "user"
    password           = "password"
    vm_console_enabled = true
  }
}

data "dcloud_vms" "test_topology_vms" {
  topology_uid = dcloud_topology.test_topology.id
}

output "vms" {
  value = data.dcloud_vms.test_topology_vms
}