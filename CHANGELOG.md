## 0.1.0 (Unreleased)

FEATURES:

* **resource/dcloud_vm**: Add `network_interfaces.assign_dhcp` to toggle a NIC's
  "Assign via DHCP" and a `dhcp_config` block (`default_gateway_ip`,
  `primary_dns_ip`, `secondary_dns_ip`) so DHCP-assigned IP, default gateway and
  DNS can be configured from Terraform instead of the UI.
* **data-source/dcloud_vms**: Expose `network_interfaces.assign_dhcp`.

BACKWARDS INCOMPATIBILITIES / NOTES:
