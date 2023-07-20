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

func dataSourceInventoryLicenses() *schema.Resource {

	return &schema.Resource{
		Description: "All the inventory licenses available to be used in a topology",

		ReadContext: dataSourceInventoryLicensesRead,

		Schema: map[string]*schema.Schema{
			"topology_uid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"inventory_licenses": {
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
						"quantity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceInventoryLicensesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	tb := m.(*tbclient.Client)

	topologyUid := d.Get("topology_uid").(string)

	inventoryLicenses, err := tb.GetAllInventoryLicenses(topologyUid)
	if err != nil {
		return diag.FromErr(err)
	}

	inventoryLicenseResources := make([]map[string]interface{}, len(inventoryLicenses))

	for i, inventoryLicense := range inventoryLicenses {
		inventoryLicenseResources[i] = convertInventoryLicenseToDataResource(inventoryLicense)
	}

	if err := d.Set("inventory_licenses", inventoryLicenseResources); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}
}

func convertInventoryLicenseToDataResource(inventoryLicense tbclient.InventoryLicense) map[string]interface{} {
	resource := make(map[string]interface{})
	resource["id"] = inventoryLicense.Id
	resource["name"] = inventoryLicense.Name
	resource["quantity"] = inventoryLicense.Quantity

	return resource
}
