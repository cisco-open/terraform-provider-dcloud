package dcloudtb

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
	"wwwin-github.cisco.com/pov-services/kapua-tb-go-client/tbclient"
)

func dataSourceHwTemplateConfigs() *schema.Resource {

	return &schema.Resource{
		Description: "All the Hardware Template Configs to be used in hardware",

		ReadContext: dataSourceHwTemplateConfigsRead,

		Schema: map[string]*schema.Schema{
			"topology_uid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"hw_template_configs": {
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

func dataSourceHwTemplateConfigsRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	tb := i.(*tbclient.Client)

	templateConfigs, err := tb.GetAllInventoryHwTemplateConfigs(data.Get("topology_uid").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	templateConfigResources := make([]map[string]interface{}, len(templateConfigs))
	for i, templateConfig := range templateConfigs {
		templateConfigResources[i] = convertHwTemplateConfigToDataResource(templateConfig)
	}

	if err := data.Set("hw_template_configs", templateConfigResources); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}
}

func convertHwTemplateConfigToDataResource(hwScript tbclient.InventoryHwScript) map[string]interface{} {
	resource := make(map[string]interface{})
	resource["uid"] = hwScript.Uid
	resource["name"] = hwScript.Name

	return resource
}
