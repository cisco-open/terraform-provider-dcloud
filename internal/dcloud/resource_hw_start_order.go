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

func resourceHwStartOrder() *schema.Resource {
	return &schema.Resource{
		Description: `Hardware Item Startup Ordering.

!>Unlike other resources there should only ever be a maximum of ONE Hardware Item Startup Order created per Topology, defining more than one will have unpredictable results.`,
		CreateContext: resourceHwStartOrderCreate,
		ReadContext:   resourceHwStartOrderRead,
		UpdateContext: resourceHwStartOrderUpdate,
		DeleteContext: resourceHwStartOrderDelete,
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
						"hw_uid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"hw_name": {
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
HW Start Order always exists within the context of a Topology, so instead of creating it we read the current one and update its values
*/
func resourceHwStartOrderCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics

	topologyUid := data.Get("topology_uid").(string)
	hwStartOrder, err := c.GetHwStartOrder(topologyUid)
	if err != nil {
		return diag.FromErr(err)
	}

	newHwStartOrder := extractHwStartOrder(data)

	hwStartOrder.Ordered = newHwStartOrder.Ordered
	hwStartOrder.Positions = newHwStartOrder.Positions

	_, err = c.UpdateHwStartOrder(*hwStartOrder)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(hwStartOrder.Uid)

	resourceHwStartOrderRead(ctx, data, i)

	return diags
}

func resourceHwStartOrderRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics

	topologyUid := data.Get("topology_uid").(string)
	hwStartOrder, err := c.GetHwStartOrder(topologyUid)
	if err != nil {
		return diag.FromErr(err)
	}

	data.Set("uid", hwStartOrder.Uid)
	data.Set("ordered", hwStartOrder.Ordered)
	data.Set("topology_uid", hwStartOrder.Topology.Uid)
	data.Set("start_positions", convertHwStartPositions(hwStartOrder.Positions))

	data.SetId(hwStartOrder.Uid)

	return diags
}

// Update and Create are the same
func resourceHwStartOrderUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	return resourceHwStartOrderCreate(ctx, data, i)
}

/*
Delete updates ordered to false, because there is always a HW Start Order present in a Topology
*/
func resourceHwStartOrderDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	topologyUid := data.Get("topology_uid").(string)
	hwStartOrder, err := c.GetHwStartOrder(topologyUid)

	hwStartOrder.Ordered = false

	_, err = c.UpdateHwStartOrder(*hwStartOrder)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func convertHwStartPositions(positions []tbclient.HwStartPosition) []map[string]interface{} {
	hwPositions := make([]map[string]interface{}, len(positions))

	sort.Slice(positions, func(i, j int) bool {
		return positions[i].Position < positions[j].Position
	})

	for i, p := range positions {
		position := make(map[string]interface{})

		position["position"] = p.Position
		position["delay_seconds"] = p.DelaySeconds
		position["hw_uid"] = p.Hw.Uid
		position["hw_name"] = p.Hw.Name

		hwPositions[i] = position
	}

	return hwPositions
}

func extractHwStartOrder(data *schema.ResourceData) *tbclient.HwStartOrder {
	return &tbclient.HwStartOrder{
		Ordered:   data.Get("ordered").(bool),
		Positions: extractHwStartPositions(data),
	}
}

func extractHwStartPositions(data *schema.ResourceData) []tbclient.HwStartPosition {
	dataPositions := data.Get("start_positions").([]interface{})

	hwStartPositions := make([]tbclient.HwStartPosition, len(dataPositions))

	for i, p := range dataPositions {
		dataPosition := p.(map[string]interface{})

		hwStartPositions[i] = tbclient.HwStartPosition{
			Position:     dataPosition["position"].(int),
			DelaySeconds: dataPosition["delay_seconds"].(int),
			Hw:           &tbclient.Hw{Uid: dataPosition["hw_uid"].(string)},
		}
	}

	return hwStartPositions
}
