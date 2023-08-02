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
)

func resourceScenario() *schema.Resource {
	return &schema.Resource{
		Description: `Scenario.

!>Unlike other resources there should only ever be a maximum of ONE Scenario resource created per Topology, defining more than one will have unpredictable results.`,
		CreateContext: resourceScenarioCreate,
		ReadContext:   resourceScenarioRead,
		UpdateContext: resourceScenarioUpdate,
		DeleteContext: resourceScenarioDelete,
		Schema: map[string]*schema.Schema{
			"uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"question": {
				Type:     schema.TypeString,
				Required: true,
			},
			"options": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internal_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Required: true,
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
Scenario always exists within the context of a Topology, so instead of creating it we read the current one and update its values
*/
func resourceScenarioCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics

	topologyUid := data.Get("topology_uid").(string)
	scenario, err := c.GetScenario(topologyUid)
	if err != nil {
		return diag.FromErr(err)
	}

	scenario.Enabled = data.Get("enabled").(bool)
	scenario.Question = data.Get("question").(string)
	scenario.Options = extractScenarioOptions(data)

	_, err = c.UpdateScenario(*scenario)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(scenario.Uid)

	resourceScenarioRead(ctx, data, i)

	return diags
}

func resourceScenarioRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics

	topologyUid := data.Get("topology_uid").(string)
	scenario, err := c.GetScenario(topologyUid)
	if err != nil {
		return diag.FromErr(err)
	}

	data.Set("uid", scenario.Uid)
	data.Set("topology_uid", scenario.Topology.Uid)
	data.Set("enabled", scenario.Enabled)
	data.Set("question", scenario.Question)
	data.Set("options", convertScenarioOptions(scenario.Options))

	data.SetId(scenario.Uid)

	return diags
}

// Update and Create are the same
func resourceScenarioUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	return resourceScenarioCreate(ctx, data, i)
}

/*
Delete updates enabled to false, because there is always a Scenario resource present in a Topology
*/
func resourceScenarioDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	topologyUid := data.Get("topology_uid").(string)
	scenario, err := c.GetScenario(topologyUid)

	scenario.Enabled = false

	_, err = c.UpdateScenario(*scenario)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func extractScenarioOptions(data *schema.ResourceData) []tbclient.ScenarioOption {
	dataOptions := data.Get("options").([]interface{})

	scenarioOptions := make([]tbclient.ScenarioOption, len(dataOptions))

	for i, p := range dataOptions {
		dataOption := p.(map[string]interface{})

		scenarioOptions[i] = tbclient.ScenarioOption{
			InternalName: dataOption["internal_name"].(string),
			DisplayName:  dataOption["display_name"].(string),
			Uid:          dataOption["uid"].(string),
		}
	}

	return scenarioOptions
}

func convertScenarioOptions(options []tbclient.ScenarioOption) []map[string]interface{} {
	scenarioOptions := make([]map[string]interface{}, len(options))

	for i, o := range options {
		option := make(map[string]interface{})

		option["uid"] = o.Uid
		option["internal_name"] = o.InternalName
		option["display_name"] = o.DisplayName

		scenarioOptions[i] = option
	}

	return scenarioOptions
}
