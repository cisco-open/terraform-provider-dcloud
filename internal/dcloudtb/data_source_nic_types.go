package dcloudtb

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
	"wwwin-github.cisco.com/pov-services/kapua-tb-go-client/tbclient"
)

func dataSourceNicTypes() *schema.Resource {

	return &schema.Resource{
		Description: "All the Nic Types available to be used in VMs",

		ReadContext: dataSourceNicTypesRead,

		Schema: map[string]*schema.Schema{
			"nic_types": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
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
		nicTypeResources[i] = converNicTypeToDataResource(nicType)
	}

	if err := data.Set("nic_types", nicTypeResources); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}
}

func converNicTypeToDataResource(nicType tbclient.NicType) map[string]interface{} {
	resource := make(map[string]interface{})
	resource["id"] = nicType.Id
	resource["name"] = nicType.Name

	return resource
}
