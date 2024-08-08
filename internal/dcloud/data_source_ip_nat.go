package dcloud

import (
	"context"
	"github.com/cisco-open/dcloud-tb-go-client/tbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func dataSourceIpNatRules() *schema.Resource {
	return &schema.Resource{
		Description: "All the IP Nat rules in a given topology",

		ReadContext: dataSourceIpNatRead,

		Schema: map[string]*schema.Schema{
			"topology_uid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_nat_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"east_west": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"scope": {
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

func dataSourceIpNatRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	tb := i.(*tbclient.Client)

	topologyUid := d.Get("topology_uid").(string)
	ipNatRules, err := tb.GetAllIpNatRules(topologyUid)

	if err != nil {
		return diag.FromErr(err)
	}

	ipNatRuleResources := make([]map[string]interface{}, len(ipNatRules))

	for i, ipNatRule := range ipNatRules {
		ipNatRuleResources[i] = convertIpNatRuleDataResource(ipNatRule)
	}
	if err := d.Set("ip_nat_rules", ipNatRuleResources); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}

}

func convertIpNatRuleDataResource(ipNatRule tbclient.IpNatRule) map[string]interface{} {
	resource := make(map[string]interface{})
	resource["uid"] = ipNatRule.Uid
	resource["topology_uid"] = ipNatRule.Topology.Uid
	resource["target_ip_address"] = ipNatRule.Target.IpAddress
	resource["target_name"] = ipNatRule.Target.Name
	resource["east_west"] = ipNatRule.EastWest
	resource["scope"] = ipNatRule.Scope

	return resource
}
