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

resource "dcloudtb_topology" "test_topology" {
  name        = "My Terraformed Topology"
  description = "A topology created from terraform"
  notes       = "Programmatic clients rule!"
  datacenter  = "LON"
}

resource "dcloudtb_network" "routed_network" {
  name                 = "A routed network"
  description          = "Demonstrating a network routed through VPOD Gateway"
  inventory_network_id = "L3-VLAN-2"
  topology_uid         = dcloudtb_topology.test_topology.id
}

resource "dcloudtb_network" "unrouted_network" {
  name                 = "An unrouted network"
  description          = "Demonstrating a network not routed through VPOD Gateway"
  inventory_network_id = "L2-VLAN-16"
  topology_uid         = dcloudtb_topology.test_topology.id
}

resource "dcloudtb_vm" "vm1" {
  inventory_vm_id   = "7668085"
  topology_uid      = dcloudtb_topology.test_topology.id
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
    network_uid = dcloudtb_network.routed_network.id
    name        = "Network adapter 0"
    mac_address = "00:50:56:00:01:AA"
    type        = "VIRTUAL_E1000"
  }

  network_interfaces {
    network_uid    = dcloudtb_network.unrouted_network.id
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

resource "dcloudtb_vm" "vm2" {
  inventory_vm_id   = "7668085"
  topology_uid      = dcloudtb_topology.test_topology.id
  name              = "Ubuntu Desktop 2"
  description       = "A standard Ubuntu Desktop VM"
  cpu_qty           = 4
  memory_mb         = 4096
  nested_hypervisor = false
  os_family         = "LINUX"

  advanced_settings {
    all_disks_non_persistent = false
    bios_uuid                = "42 3a 5f 9d f1 a8 7c 0e-7d c2 44 27 2e d6 67 ad"
    name_in_hypervisor       = "ubuntu"
    not_started              = false
  }

  remote_access {
    vm_console_enabled = false
  }

  network_interfaces {
    network_uid = dcloudtb_network.routed_network.id
    name        = "Network adapter 0"
    mac_address = "00:50:56:00:01:AF"
    type        = "VIRTUAL_E1000"
  }
}


resource "dcloudtb_hw" "hw1" {
  topology_uid               = dcloudtb_topology.test_topology.id
  inventory_hw_id            = "76"
  name                       = "IE 4000 Device"
  hardware_console_enabled   = false
  startup_script_uid         = "bjlfkxev55nh35eh6kku13971"
  custom_script_uid          = "668eljku7jwpk8bpysz5njyrz"
  shutdown_script_uid        = "435ya6tjh5u4uv3ku2kphesr"
  template_config_script_uid = "79ila00mn7icfbtk3dg7fuasy"

  network_interfaces {
    network_interface_id = "GigabitEthernet1/0/24"
    network_uid          = dcloudtb_network.routed_network.id
  }
}

resource "dcloudtb_hw" "hw2" {
  topology_uid    = dcloudtb_topology.test_topology.id
  inventory_hw_id = "14"
  name            = "UCS Hardware Pod"
}