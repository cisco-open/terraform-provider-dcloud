terraform {
  required_providers {
    dcloud = {
      version = "0.1"
      source  = "cisco.com/cisco-open/dcloud"
    }
  }
}

provider "dcloud" {
  tb_url = "https://tbv3-production.ciscodcloud.com/api"
}

resource "dcloud_topology" "test_topology" {
  name        = "My Terraformed Topology"
  description = "A topology created from terraform"
  notes       = "Programmatic clients rule!"
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

resource "dcloud_vm" "vm2" {
  inventory_vm_id   = "7668085"
  topology_uid      = dcloud_topology.test_topology.id
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
    network_uid = dcloud_network.routed_network.id
    name        = "Network adapter 0"
    mac_address = "00:50:56:00:01:AF"
    type        = "VIRTUAL_E1000"
  }
}
/*

resource "dcloud_hw" "hw1" {
  topology_uid               = dcloud_topology.test_topology.id
  inventory_hw_id            = "76"
  name                       = "IE 4000 Device"
  hardware_console_enabled   = false
  startup_script_uid         = "bjlfkxev55nh35eh6kku13971"
  custom_script_uid          = "668eljku7jwpk8bpysz5njyrz"
  shutdown_script_uid        = "435ya6tjh5u4uv3ku2kphesr"
  template_config_script_uid = "79ila00mn7icfbtk3dg7fuasy"

  network_interfaces {
    network_interface_id = "GigabitEthernet1/0/24"
    network_uid          = dcloud_network.routed_network.id
  }
}

resource "dcloud_hw" "hw2" {
  topology_uid    = dcloud_topology.test_topology.id
  inventory_hw_id = "14"
  name            = "UCS Hardware Pod"
}

resource "dcloud_license" "mc_license" {
  quantity             = 3
  inventory_license_id = "340"
  topology_uid         = dcloud_topology.test_topology.id
}

resource "dcloud_vm_start_order" "vm_start_order" {
  ordered      = true
  topology_uid = dcloud_topology.test_topology.id

  start_positions {
    position      = 1
    delay_seconds = 10
    vm_uid        = dcloud_vm.vm1.id
  }

  start_positions {
    position      = 2
    delay_seconds = 20
    vm_uid        = dcloud_vm.vm2.id
  }
}

resource "dcloud_vm_stop_order" "vm_stop_order" {
  ordered      = true
  topology_uid = dcloud_topology.test_topology.id

  stop_positions {
    position = 1
    vm_uid   = dcloud_vm.vm2.id
  }

  stop_positions {
    position = 2
    vm_uid   = dcloud_vm.vm1.id
  }
}

resource "dcloud_hw_start_order" "hw_start_order" {
  ordered      = true
  topology_uid = dcloud_topology.test_topology.id

  start_positions {
    position      = 1
    delay_seconds = 10
    hw_uid        = dcloud_hw.hw1.id
  }

  start_positions {
    position      = 2
    delay_seconds = 20
    hw_uid        = dcloud_hw.hw2.id
  }
}

resource "dcloud_remote_access" "remote_access" {
  any_connect_enabled  = true
  endpoint_kit_enabled = true
  topology_uid         = dcloud_topology.test_topology.id
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

resource "dcloud_documentation" "documentation" {
  topology_uid = dcloud_topology.test_topology.id
  doc_url = "https://johndoe.com"
}

resource "dcloud_telephony" "telephony" {
  topology_uid = dcloud_topology.test_topology.id
  inventory_telephony_id = "1"
}
*/

resource "dcloud_vm" "vm3" {
  inventory_vm_id   = "7668085"
  topology_uid      = dcloud_topology.test_topology.id
  name              = "Ubuntu Desktop 3"
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
    ip_address     = "127.0.0.3"
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


resource "dcloud_ip_nat_rule" "ip_nat_rule"{
  topology_uid = dcloud_topology.test_topology.id
  target_ip_address = "192.168.1.1"
  target_name = "Sample Device"
  east_west = false
}

resource "dcloud_inbound_proxy_rule" "inbound_proxy_rule"{
  topology_uid = dcloud_topology.test_topology.id
  nic_uid = dcloud_vm.vm1.network_interfaces[1].uid
  target_vm_name = dcloud_vm.vm1.name
  tcp_port = 443
  url_path = "/demo/demo010/"
  hyperlink = "Click Me"
  ssl = true
}
resource "dcloud_inbound_proxy_rule" "inbound_proxy_rule_2"{
  topology_uid = dcloud_topology.test_topology.id
  nic_uid = dcloud_vm.vm3.network_interfaces[1].uid
  target_vm_name = dcloud_vm.vm1.name
  tcp_port = 443
  url_path = "/demo/demo010/"
  hyperlink = "Click Me"
  ssl = true
}

data "dcloud_inbound_proxy_rules" "test_topology_inbound_proxy" {
  topology_uid = dcloud_topology.test_topology.id
}

output "inbound_proxy_rules" {
  value = data.dcloud_inbound_proxy_rules.test_topology_inbound_proxy
}
