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

func resourceLicense() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLicenseCreate,
		ReadContext:   resourceLicenseRead,
		UpdateContext: resourceLicenseUpdate,
		DeleteContext: resourceLicenseDelete,
		Schema: map[string]*schema.Schema{
			"uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"quantity": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"inventory_license_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"inventory_license_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"topology_uid": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceLicenseCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics

	license := tbclient.License{
		Quantity: data.Get("quantity").(int),
		InventoryLicense: &tbclient.InventoryLicense{
			Id: data.Get("inventory_license_id").(string),
		},
		Topology: &tbclient.Topology{
			Uid: data.Get("topology_uid").(string),
		},
	}

	l, err := c.CreateLicense(license)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(l.Uid)

	resourceLicenseRead(ctx, data, i)

	return diags
}

func resourceLicenseRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics

	license, err := c.GetLicense(data.Id())
	if err != nil {
		return handleClientError(err, data, diags)
	}

	data.Set("uid", license.Uid)
	data.Set("quantity", license.Quantity)
	data.Set("inventory_license_id", license.InventoryLicense.Id)
	data.Set("inventory_license_name", license.InventoryLicense.Name)
	data.Set("topology_uid", license.Topology.Uid)

	return diags
}

func resourceLicenseUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	license := tbclient.License{
		Uid:      data.Get("uid").(string),
		Quantity: data.Get("quantity").(int),
		InventoryLicense: &tbclient.InventoryLicense{
			Id: data.Get("inventory_license_id").(string),
		},
		Topology: &tbclient.Topology{
			Uid: data.Get("topology_uid").(string),
		},
	}

	_, err := c.UpdateLicense(license)
	if err != nil {
		var diags diag.Diagnostics
		return handleClientError(err, data, diags)
	}

	return resourceLicenseRead(ctx, data, i)
}

func resourceLicenseDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	if err := c.DeleteLicense(data.Id()); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
