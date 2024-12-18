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
	"strconv"
	"time"
)

func dataSourceEvcModes() *schema.Resource {

	return &schema.Resource{
		Description: "All the EVC Modes available to be used with VMs",

		ReadContext: dataSourceEvcModesRead,

		Schema: map[string]*schema.Schema{
			"evc_modes": {
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
					},
				},
			},
		},
	}
}

func dataSourceEvcModesRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	tb := m.(*tbclient.Client)

	evcModes, err := tb.GetAllEvcModes()
	if err != nil {
		return diag.FromErr(err)
	}

	evcModeResources := make([]map[string]interface{}, len(evcModes))
	for i, evcMode := range evcModes {
		evcModeResources[i] = convertEvcModeToDataResource(evcMode)
	}

	if err := data.Set("evc_modes", evcModeResources); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}
}

func convertEvcModeToDataResource(evcMode tbclient.EvcMode) map[string]interface{} {
	resource := make(map[string]interface{})
	resource["id"] = evcMode.Id
	resource["name"] = evcMode.Name

	return resource
}
