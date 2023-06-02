// Copyright 2023 Cisco Systems, Inc. and its affiliates
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package dcloudtb

import (
	"context"
	"github.com/cisco-open/kapua-tb-go-client/tbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"sort"
	"strconv"
	"time"
)

func dataSourceVms() *schema.Resource {

	return &schema.Resource{
		Description: "All the VMs currently in a given topology",

		ReadContext: dataSourceVmsRead,

		Schema: map[string]*schema.Schema{
			"topology_uid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory_mb": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cpu_qty": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"nested_hypervisor": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"inventory_vm_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_family": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"advanced_settings": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name_in_hypervisor": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"bios_uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"not_started": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"all_disks_non_persistent": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"remote_access": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"password": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vm_console_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"display_credentials": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"username": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"password": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"internal_urls": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"location": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"description": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"guest_automation": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"command": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"delay_seconds": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"network_interfaces": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"uid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"mac_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"network_uid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ssh_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"rdp_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"rdp_auto_login": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"topology_uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceVmsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	tb := m.(*tbclient.Client)

	topologyUid := d.Get("topology_uid").(string)

	vms, err := tb.GetAllVms(topologyUid)
	if err != nil {
		return diag.FromErr(err)
	}

	vmResources := make([]map[string]interface{}, len(vms))

	for i, vm := range vms {
		vmResources[i] = convertVmDataResource(vm, topologyUid)
	}

	if err := d.Set("vms", vmResources); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}
}

func convertVmDataResource(vm tbclient.Vm, topologyUid string) map[string]interface{} {
	resource := make(map[string]interface{})

	addVmFields(vm, resource)
	if as := convertAdvancedSettings(vm); as != nil {
		resource["advanced_settings"] = as
	}
	if ra := convertRemoteAccess(vm); ra != nil {
		resource["remote_access"] = ra
	}
	if ga := convertGuestAutomation(vm); ga != nil {
		resource["guest_automation"] = ga
	}
	if nics := convertNics(vm); nics != nil {
		resource["network_interfaces"] = nics
	}

	return resource
}

func addVmFields(vm tbclient.Vm, resource map[string]interface{}) {
	resource["uid"] = vm.Uid
	resource["name"] = vm.Name
	resource["description"] = vm.Description
	resource["cpu_qty"] = vm.CpuQty
	resource["memory_mb"] = vm.MemoryMb
	resource["nested_hypervisor"] = vm.NestedHypervisor
	resource["inventory_vm_id"] = vm.InventoryVmId
	resource["os_family"] = vm.OsFamily
	resource["topology_uid"] = vm.Topology.Uid
}

func convertAdvancedSettings(vm tbclient.Vm) []interface{} {
	if advancedSettings := vm.AdvancedSettings; advancedSettings != nil {
		as := make(map[string]interface{})

		as["name_in_hypervisor"] = advancedSettings.NameInHypervisor
		as["bios_uuid"] = advancedSettings.BiosUuid
		as["not_started"] = advancedSettings.NotStarted
		as["all_disks_non_persistent"] = advancedSettings.AllDisksNonPersistent

		return []interface{}{as}
	}
	return nil
}

func convertRemoteAccess(vm tbclient.Vm) []interface{} {
	if remoteAccess := vm.RemoteAccess; remoteAccess != nil {
		ra := make(map[string]interface{})

		ra["username"] = remoteAccess.Username
		ra["password"] = remoteAccess.Password
		ra["vm_console_enabled"] = remoteAccess.VmConsoleEnabled

		addDisplayCredentials(remoteAccess, ra)
		addInternalUrls(remoteAccess, ra)

		return []interface{}{ra}
	}
	return nil
}

func addDisplayCredentials(remoteAccess *tbclient.VmRemoteAccess, ra map[string]interface{}) {
	if displayCredentials := remoteAccess.DisplayCredentials; displayCredentials != nil {
		dc := make(map[string]interface{})

		dc["username"] = displayCredentials.Username
		dc["password"] = displayCredentials.Password

		ra["display_credentials"] = []interface{}{dc}
	}
}

func addInternalUrls(remoteAccess *tbclient.VmRemoteAccess, ra map[string]interface{}) {
	internalUrls := make([]map[string]interface{}, len(remoteAccess.InternalUrls))
	for i, internalUrl := range remoteAccess.InternalUrls {
		iu := make(map[string]interface{})

		iu["location"] = internalUrl.Location
		iu["description"] = internalUrl.Description

		internalUrls[i] = iu
	}
	ra["internal_urls"] = internalUrls
}

func convertGuestAutomation(vm tbclient.Vm) []interface{} {
	if guestAutomation := vm.GuestAutomation; guestAutomation != nil {
		ga := make(map[string]interface{})

		ga["command"] = guestAutomation.Command
		ga["delay_seconds"] = guestAutomation.DelaySecs

		return []interface{}{ga}
	}
	return nil
}

func convertNics(vm tbclient.Vm) []map[string]interface{} {
	nics := make([]map[string]interface{}, len(vm.VmNetworkInterfaces))

	sort.Slice(vm.VmNetworkInterfaces, func(i, j int) bool {
		return vm.VmNetworkInterfaces[i].Name < vm.VmNetworkInterfaces[j].Name
	})

	for i, nic := range vm.VmNetworkInterfaces {
		n := make(map[string]interface{})

		n["uid"] = nic.Uid
		n["name"] = nic.Name
		n["mac_address"] = nic.MacAddress
		n["ip_address"] = nic.IpAddress
		n["type"] = nic.Type
		n["network_uid"] = nic.Network.Uid

		if rdp := nic.Rdp; rdp != nil {
			n["rdp_enabled"] = rdp.Enabled
			n["rdp_auto_login"] = rdp.AutoLogin
		} else {
			n["rdp_enabled"] = false
			n["rdp_auto_login"] = false
		}

		if ssh := nic.Ssh; ssh != nil {
			n["ssh_enabled"] = ssh.Enabled
		} else {
			n["ssh_enabled"] = false
		}

		nics[i] = n
	}
	return nics
}
