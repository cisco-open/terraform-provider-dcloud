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

func dataSourceNicTypes() *schema.Resource {

	return &schema.Resource{
		Description: "All the Nic Types available to be used in VMs",

		ReadContext: dataSourceNicTypesRead,

		Schema: map[string]*schema.Schema{
			"nic_types": {
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

func dataSourceNicTypesRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	tb := m.(*tbclient.Client)

	nicTypes, err := tb.GetAllNicTypes()
	if err != nil {
		return diag.FromErr(err)
	}

	nicTypeResources := make([]map[string]interface{}, len(nicTypes))
	for i, nicType := range nicTypes {
		nicTypeResources[i] = convertNicTypeToDataResource(nicType)
	}

	if err := data.Set("nic_types", nicTypeResources); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}
}

func convertNicTypeToDataResource(nicType tbclient.NicType) map[string]interface{} {
	resource := make(map[string]interface{})
	resource["id"] = nicType.Id
	resource["name"] = nicType.Name

	return resource
}
