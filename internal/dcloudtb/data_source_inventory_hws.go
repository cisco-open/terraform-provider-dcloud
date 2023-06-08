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
	"github.com/cisco-open/dcloud-tb-go-client/tbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func dataSourceInventoryHws() *schema.Resource {

	return &schema.Resource{
		Description: "All the inventory HW Items available to be used in a topology",

		ReadContext: dataSourceInventoryHwsRead,

		Schema: map[string]*schema.Schema{
			"topology_uid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"inventory_hws": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
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
						"power_control_available": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"hardware_console_available": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"network_interfaces": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceInventoryHwsRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {

	tb := i.(*tbclient.Client)

	topologyUid := d.Get("topology_uid").(string)

	inventoryHws, err := tb.GetAllInventoryHws(topologyUid)
	if err != nil {
		return diag.FromErr(err)
	}

	inventoryHwResources := make([]map[string]interface{}, len(inventoryHws))

	for i, inventoryHw := range inventoryHws {
		inventoryHwResources[i] = convertInventoryHwToDataResource(inventoryHw, topologyUid)
	}

	if err := d.Set("inventory_hws", inventoryHwResources); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}
}

func convertInventoryHwToDataResource(inventoryHw tbclient.InventoryHw, topologyUid string) map[string]interface{} {
	resource := make(map[string]interface{})
	resource["id"] = inventoryHw.Id
	resource["name"] = inventoryHw.Name
	resource["description"] = inventoryHw.Description
	resource["power_control_available"] = inventoryHw.PowerControlAvailable
	resource["hardware_console_available"] = inventoryHw.HardwareConsoleAvailable

	nics := make([]interface{}, len(inventoryHw.NetworkInterfaces))

	for i, nic := range inventoryHw.NetworkInterfaces {
		nicResource := make(map[string]interface{})
		nicResource["id"] = nic.Id

		nics[i] = nicResource
	}
	resource["network_interfaces"] = nics

	return resource
}
