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

func dataSourceHwScripts() *schema.Resource {

	return &schema.Resource{
		Description: "All the Hardware Scripts to be used in hardware",

		ReadContext: dataSourceHwScriptsRead,

		Schema: map[string]*schema.Schema{
			"topology_uid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"hw_scripts": {
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
					},
				},
			},
		},
	}
}

func dataSourceHwScriptsRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	tb := i.(*tbclient.Client)

	hwScripts, err := tb.GetAllInventoryHwScripts(data.Get("topology_uid").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	hwScriptResources := make([]map[string]interface{}, len(hwScripts))
	for i, hwScript := range hwScripts {
		hwScriptResources[i] = convertHwScriptToDataResource(hwScript)
	}

	if err := data.Set("hw_scripts", hwScriptResources); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}
}

func convertHwScriptToDataResource(hwScript tbclient.InventoryHwScript) map[string]interface{} {
	resource := make(map[string]interface{})
	resource["uid"] = hwScript.Uid
	resource["name"] = hwScript.Name

	return resource
}
