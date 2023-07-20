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

func resourceVmStartOrder() *schema.Resource {
	return &schema.Resource{
		Description: `Virtual Machine Startup Ordering.

!>Unlike other resources there should only ever be a maximum of ONE Virtual Machine Startup Order created per Topology, defining more than one will have unpredictable results.`,
		CreateContext: resourceVmStartOrderCreate,
		ReadContext:   resourceVmStartOrderRead,
		UpdateContext: resourceVmStartOrderUpdate,
		DeleteContext: resourceVmStartOrderDelete,
		Schema: map[string]*schema.Schema{
			"uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ordered": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"start_positions": {
				Description: "\nThe ordered collection of Start Positions, always specify in ascending position order within your .tf file to avoid unnecessary churn on apply.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"position": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"delay_seconds": {
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
VM Start Order always exists within the context of a Topology, so instead of creating it we read the current one and update its values
*/
func resourceVmStartOrderCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics

	topologyUid := data.Get("topology_uid").(string)
	vmStartOrder, err := c.GetVmStartOrder(topologyUid)
	if err != nil {
		return diag.FromErr(err)
	}

	newVmStartOrder := extractVmStartOrder(data)

	vmStartOrder.Ordered = newVmStartOrder.Ordered
	vmStartOrder.Positions = newVmStartOrder.Positions

	_, err = c.UpdateVmStartOrder(*vmStartOrder)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(vmStartOrder.Uid)

	resourceVmStartOrderRead(ctx, data, i)

	return diags
}

func resourceVmStartOrderRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics

	topologyUid := data.Get("topology_uid").(string)
	vmStartOrder, err := c.GetVmStartOrder(topologyUid)
	if err != nil {
		return diag.FromErr(err)
	}

	data.Set("uid", vmStartOrder.Uid)
	data.Set("ordered", vmStartOrder.Ordered)
	data.Set("topology_uid", vmStartOrder.Topology.Uid)
	data.Set("start_positions", convertStartPositions(vmStartOrder.Positions))

	data.SetId(vmStartOrder.Uid)

	return diags
}

// Update and Create are the same
func resourceVmStartOrderUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	return resourceVmStartOrderCreate(ctx, data, i)
}

/*
Delete updates ordered to false, because there is always a VM Start Order present in a Topology
*/
func resourceVmStartOrderDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	topologyUid := data.Get("topology_uid").(string)
	vmStartOrder, err := c.GetVmStartOrder(topologyUid)

	vmStartOrder.Ordered = false

	_, err = c.UpdateVmStartOrder(*vmStartOrder)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func convertStartPositions(positions []tbclient.VmStartPosition) []map[string]interface{} {
	vmPositions := make([]map[string]interface{}, len(positions))

	sort.Slice(positions, func(i, j int) bool {
		return positions[i].Position < positions[j].Position
	})

	for i, p := range positions {
		position := make(map[string]interface{})

		position["position"] = p.Position
		position["delay_seconds"] = p.DelaySeconds
		position["vm_uid"] = p.Vm.Uid
		position["vm_name"] = p.Vm.Name

		vmPositions[i] = position
	}

	return vmPositions
}

func extractVmStartOrder(data *schema.ResourceData) *tbclient.VmStartOrder {
	return &tbclient.VmStartOrder{
		Ordered:   data.Get("ordered").(bool),
		Positions: extractPositions(data),
	}
}

func extractPositions(data *schema.ResourceData) []tbclient.VmStartPosition {
	dataPositions := data.Get("start_positions").([]interface{})

	vmStartPositions := make([]tbclient.VmStartPosition, len(dataPositions))

	for i, p := range dataPositions {
		dataPosition := p.(map[string]interface{})

		vmStartPositions[i] = tbclient.VmStartPosition{
			Position:     dataPosition["position"].(int),
			DelaySeconds: dataPosition["delay_seconds"].(int),
			Vm:           &tbclient.Vm{Uid: dataPosition["vm_uid"].(string)},
		}
	}

	return vmStartPositions
}
