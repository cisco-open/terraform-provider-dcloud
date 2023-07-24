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

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TB_AUTH_TOKEN", nil),
				Sensitive:   true,
			},
			"tb_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"debug": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"dcloud_topologies":                    dataSourceTopologies(),
			"dcloud_inventory_networks":            dataSourceInventoryNetworks(),
			"dcloud_networks":                      dataSourceNetworks(),
			"dcloud_os_families":                   dataSourceOsFamilies(),
			"dcloud_nic_types":                     dataSourceNicTypes(),
			"dcloud_inventory_vms":                 dataSourceInventoryVms(),
			"dcloud_vms":                           dataSourceVms(),
			"dcloud_inventory_hw_scripts":          dataSourceHwScripts(),
			"dcloud_inventory_hw_template_configs": dataSourceHwTemplateConfigs(),
			"dcloud_inventory_hws":                 dataSourceInventoryHws(),
			"dcloud_hws":                           dataSourceHws(),
			"dcloud_inventory_licenses":            dataSourceInventoryLicenses(),
			"dcloud_licenses":                      dataSourceLicenses(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"dcloud_topology":       resourceTopology(),
			"dcloud_network":        resourceNetwork(),
			"dcloud_vm":             resourceVm(),
			"dcloud_hw":             resourceHw(),
			"dcloud_license":        resourceLicense(),
			"dcloud_vm_start_order": resourceVmStartOrder(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
	authToken := data.Get("auth_token").(string)
	url := data.Get("tb_url").(string)
	debug := data.Get("debug").(bool)

	var diags diag.Diagnostics

	c := tbclient.NewClient(&url, &authToken)
	c.Debug = debug
	c.UserAgent = "terraform-provider-dcloud/v0.0.1" // TODO - replace with actual application version, if possible

	return c, diags
}
