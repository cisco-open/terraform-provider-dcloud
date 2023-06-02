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
	"strconv"
	"time"
)

func dataSourceOsFamilies() *schema.Resource {

	return &schema.Resource{
		Description: "All the OS Families available to be used in VMs",

		ReadContext: dataSourceOsFamiliesRead,

		Schema: map[string]*schema.Schema{
			"os_families": {
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

func dataSourceOsFamiliesRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	tb := m.(*tbclient.Client)

	osFamilies, err := tb.GetAllOsFamilies()
	if err != nil {
		return diag.FromErr(err)
	}

	osFamilyResources := make([]map[string]interface{}, len(osFamilies))
	for i, osFamily := range osFamilies {
		osFamilyResources[i] = convertOsFamilyToDataResource(osFamily)
	}

	if err := data.Set("os_families", osFamilyResources); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}
}

func convertOsFamilyToDataResource(osFamily tbclient.OsFamily) map[string]interface{} {
	resource := make(map[string]interface{})
	resource["id"] = osFamily.Id
	resource["name"] = osFamily.Name

	return resource
}
