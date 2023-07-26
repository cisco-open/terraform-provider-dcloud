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

func resourceRemoteAccess() *schema.Resource {
	return &schema.Resource{
		Description: `Remote Access.

!>Unlike other resources there should only ever be a maximum of ONE Remote Access resource created per Topology, defining more than one will have unpredictable results.`,
		CreateContext: resourceRemoteAccessCreate,
		ReadContext:   resourceRemoteAccessRead,
		UpdateContext: resourceRemoteAccessUpdate,
		DeleteContext: resourceRemoteAccessDelete,
		Schema: map[string]*schema.Schema{
			"uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"any_connect_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"endpoint_kit_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"topology_uid": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

/*
Remote Access always exists within the context of a Topology, so instead of creating it we read the current one and update its values
*/
func resourceRemoteAccessCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics

	topologyUid := data.Get("topology_uid").(string)
	remoteAccess, err := c.GetRemoteAccess(topologyUid)
	if err != nil {
		return diag.FromErr(err)
	}

	remoteAccess.AnyconnectEnabled = data.Get("any_connect_enabled").(bool)
	remoteAccess.EndpointKitEnabled = data.Get("endpoint_kit_enabled").(bool)

	_, err = c.UpdateRemoteAccess(*remoteAccess)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(remoteAccess.Uid)

	resourceRemoteAccessRead(ctx, data, i)

	return diags
}

func resourceRemoteAccessRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics

	topologyUid := data.Get("topology_uid").(string)
	remoteAccess, err := c.GetRemoteAccess(topologyUid)
	if err != nil {
		return diag.FromErr(err)
	}

	data.Set("uid", remoteAccess.Uid)
	data.Set("topology_uid", remoteAccess.Topology.Uid)
	data.Set("any_connect_enabled", remoteAccess.AnyconnectEnabled)
	data.Set("endpoint_kit_enabled", remoteAccess.EndpointKitEnabled)

	data.SetId(remoteAccess.Uid)

	return diags
}

// Update and Create are the same
func resourceRemoteAccessUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	return resourceRemoteAccessCreate(ctx, data, i)
}

/*
Delete updates any connect and endpoint kit to false, because there is always a Remote Access resource present in a Topology
*/
func resourceRemoteAccessDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	topologyUid := data.Get("topology_uid").(string)
	remoteAccess, err := c.GetRemoteAccess(topologyUid)

	remoteAccess.AnyconnectEnabled = false
	remoteAccess.EndpointKitEnabled = false

	_, err = c.UpdateRemoteAccess(*remoteAccess)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
