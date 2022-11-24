package dcloudtb

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
	"wwwin-github.cisco.com/pov-services/kapua-tb-go-client/tbclient"
)

func dataSourceTopologies() *schema.Resource {

	return &schema.Resource{
		Description: "All the topologies owned or shared to the authenticated user",

		ReadContext: dataSourceTopologiesRead,

		Schema: map[string]*schema.Schema{
			"topologies": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: topologyDataResourceSchema,
				},
			},
		},
	}
}

func dataSourceTopologiesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	tb := m.(*tbclient.Client)

	topologies, err := tb.GetAllTopologies()
	if err != nil {
		return diag.FromErr(err)
	}

	topologyResources := make([]map[string]interface{}, len(topologies))

	for i, topology := range topologies {
		topologyResources[i] = convert(topology)
	}

	if err := d.Set("topologies", topologyResources); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}
}
