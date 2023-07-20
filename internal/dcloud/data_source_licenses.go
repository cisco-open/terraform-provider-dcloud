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

func dataSourceLicenses() *schema.Resource {

	return &schema.Resource{
		Description: "All the licenses currently in a given topology",

		ReadContext: dataSourceLicenseRead,

		Schema: map[string]*schema.Schema{
			"topology_uid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"licenses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quantity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"inventory_license_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"inventory_license_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"topology_uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceLicenseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	tb := m.(*tbclient.Client)

	topologyUid := d.Get("topology_uid").(string)

	licenses, err := tb.GetAllLicenses(topologyUid)
	if err != nil {
		return diag.FromErr(err)
	}

	licenseResources := make([]map[string]interface{}, len(licenses))

	for i, license := range licenses {
		licenseResources[i] = convertLicenseToDataResource(license)
	}

	if err := d.Set("licenses", licenseResources); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}
}

func convertLicenseToDataResource(license tbclient.License) map[string]interface{} {
	resource := make(map[string]interface{})
	resource["uid"] = license.Uid
	resource["quantity"] = license.Quantity
	resource["inventory_license_id"] = license.InventoryLicense.Id
	resource["inventory_license_name"] = license.InventoryLicense.Name
	resource["topology_uid"] = license.Topology.Uid

	return resource
}
