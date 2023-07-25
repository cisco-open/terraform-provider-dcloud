// Copyright 2023 Cisco Systems, Inc. and its affiliates
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package dcloud

import (
	"context"
	"github.com/cisco-open/dcloud-tb-go-client/tbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"sort"
)

func resourceVmStopOrder() *schema.Resource {
	return &schema.Resource{
		Description: `Virtual Machine Stop Ordering.

!>Unlike other resources there should only ever be a maximum of ONE Virtual Machine Stop Order created per Topology, defining more than one will have unpredictable results.`,
		CreateContext: resourceVmStopOrderCreate,
		ReadContext:   resourceVmStopOrderRead,
		UpdateContext: resourceVmStopOrderUpdate,
		DeleteContext: resourceVmStopOrderDelete,
		Schema: map[string]*schema.Schema{
			"uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ordered": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"stop_positions": {
				Description: "\nThe ordered collection of Stop Positions, always specify in ascending position order within your .tf file to avoid unnecessary churn on apply.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"position": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"vm_uid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"vm_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"topology_uid": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

/*
VM Stop Order always exists within the context of a Topology, so instead of creating it we read the current one and update its values
*/
func resourceVmStopOrderCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics

	topologyUid := data.Get("topology_uid").(string)
	vmStopOrder, err := c.GetVmStopOrder(topologyUid)
	if err != nil {
		return diag.FromErr(err)
	}

	newVmStopOrder := extractVmStopOrder(data)

	vmStopOrder.Ordered = newVmStopOrder.Ordered
	vmStopOrder.Positions = newVmStopOrder.Positions

	_, err = c.UpdateVmStopOrder(*vmStopOrder)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(vmStopOrder.Uid)

	resourceVmStopOrderRead(ctx, data, i)

	return diags
}

func resourceVmStopOrderRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics

	topologyUid := data.Get("topology_uid").(string)
	vmStopOrder, err := c.GetVmStopOrder(topologyUid)
	if err != nil {
		return diag.FromErr(err)
	}

	data.Set("uid", vmStopOrder.Uid)
	data.Set("ordered", vmStopOrder.Ordered)
	data.Set("topology_uid", vmStopOrder.Topology.Uid)
	data.Set("stop_positions", convertVmStopPositions(vmStopOrder.Positions))

	data.SetId(vmStopOrder.Uid)

	return diags
}

// Update and Create are the same
func resourceVmStopOrderUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	return resourceVmStopOrderCreate(ctx, data, i)
}

/*
Delete updates ordered to false, because there is always a VM Stop Order present in a Topology
*/
func resourceVmStopOrderDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	topologyUid := data.Get("topology_uid").(string)
	vmStopOrder, err := c.GetVmStopOrder(topologyUid)

	vmStopOrder.Ordered = false

	_, err = c.UpdateVmStopOrder(*vmStopOrder)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func extractVmStopOrder(data *schema.ResourceData) *tbclient.VmStopOrder {
	return &tbclient.VmStopOrder{
		Ordered:   data.Get("ordered").(bool),
		Positions: extractVmStopPositions(data),
	}
}

func extractVmStopPositions(data *schema.ResourceData) []tbclient.VmStopPosition {
	dataPositions := data.Get("stop_positions").([]interface{})

	vmStopPositions := make([]tbclient.VmStopPosition, len(dataPositions))

	for i, p := range dataPositions {
		dataPosition := p.(map[string]interface{})

		vmStopPositions[i] = tbclient.VmStopPosition{
			Position: dataPosition["position"].(int),
			Vm:       &tbclient.Vm{Uid: dataPosition["vm_uid"].(string)},
		}
	}

	return vmStopPositions
}

func convertVmStopPositions(positions []tbclient.VmStopPosition) []map[string]interface{} {
	vmPositions := make([]map[string]interface{}, len(positions))

	sort.Slice(positions, func(i, j int) bool {
		return positions[i].Position < positions[j].Position
	})

	for i, p := range positions {
		position := make(map[string]interface{})

		position["position"] = p.Position
		position["vm_uid"] = p.Vm.Uid
		position["vm_name"] = p.Vm.Name

		vmPositions[i] = position
	}

	return vmPositions
}
