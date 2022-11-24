package dcloudtb

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"wwwin-github.cisco.com/pov-services/kapua-tb-go-client/tbclient"
)

var topologyDataResourceSchema = map[string]*schema.Schema{
	"uid": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	},
	"name": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	},
	"description": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	},
	"datacenter": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	},
	"notes": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	},
	"status": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	},
}

var topologyResourceSchema = map[string]*schema.Schema{
	"uid": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	},
	"name": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	},
	"description": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	},
	"datacenter": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	},
	"notes": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	},
	"status": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	},
}

func convert(topology tbclient.Topology) map[string]interface{} {
	resource := make(map[string]interface{})
	resource["uid"] = topology.Uid
	resource["name"] = topology.Name
	resource["description"] = topology.Description
	resource["datacenter"] = topology.Datacenter
	resource["notes"] = topology.Notes
	resource["status"] = topology.Status

	return resource
}
