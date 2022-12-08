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
			"topologies": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datacenter": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notes": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
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
		topologyResources[i] = convertTopologyToDataResource(topology)
	}

	if err := d.Set("topologies", topologyResources); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}
}

func convertTopologyToDataResource(topology tbclient.Topology) map[string]interface{} {
	resource := make(map[string]interface{})
	resource["uid"] = topology.Uid
	resource["name"] = topology.Name
	resource["description"] = topology.Description
	resource["datacenter"] = topology.Datacenter
	resource["notes"] = topology.Notes
	resource["status"] = topology.Status

	return resource
}
